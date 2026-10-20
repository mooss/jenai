package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type siblingSemantics string

const (
	siblingError  siblingSemantics = "error"
	siblingIgnore siblingSemantics = "ignore"
	siblingMerge  siblingSemantics = "merge"
	siblingAllOf  siblingSemantics = "allOf"
)

type resolver struct {
	docs      map[string]any
	client    *http.Client
	siblings  siblingSemantics
}

func documentURL(u *url.URL) *url.URL {
	v := *u
	v.Fragment = ""
	v.RawFragment = ""
	return &v
}

func inputURL(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err == nil && (u.Scheme == "http" ||
		u.Scheme == "https" || u.Scheme == "file") {
		return u, nil
	}

	path, err := filepath.Abs(s)
	if err != nil {
		return nil, err
	}
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return &url.URL{Scheme: "file", Path: path}, nil
}

func (r *resolver) load(u *url.URL) (any, error) {
	u = documentURL(u)
	key := u.String()
	if doc, ok := r.docs[key]; ok {
		return doc, nil
	}

	var reader io.ReadCloser
	switch u.Scheme {
	case "file":
		if u.Host != "" && u.Host != "localhost" {
			return nil, fmt.Errorf("unsupported file host: %s", u.Host)
		}
		path := u.Path
		if os.PathSeparator == '\\' &&
			len(path) >= 3 && path[0] == '/' && path[2] == ':' {
			path = path[1:]
		}
		f, err := os.Open(filepath.FromSlash(path))
		if err != nil {
			return nil, err
		}
		reader = f

	case "http", "https":
		resp, err := r.client.Get(key)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return nil, fmt.Errorf("%s: %s", key, resp.Status)
		}
		reader = resp.Body

	default:
		return nil, fmt.Errorf("unsupported URL scheme: %q", u.Scheme)
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)
	dec.UseNumber()

	var doc any
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	// Require exactly one JSON value.
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("%s: multiple JSON values", key)
		}
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	r.docs[key] = doc
	return doc, nil
}

// u.Fragment has already been URI-percent-decoded by net/url.
func pointer(doc any, fragment string) (any, error) {
	if fragment == "" {
		return doc, nil
	}
	if !strings.HasPrefix(fragment, "/") {
		return nil, fmt.Errorf(
			"unsupported fragment %q: expected a JSON Pointer", fragment)
	}

	value := doc
	for raw := range strings.SplitSeq(fragment[1:], "/") {
		// Validate JSON Pointer escape sequences.
		for i := 0; i < len(raw); i++ {
			if raw[i] == '~' {
				if i+1 == len(raw) ||
					(raw[i+1] != '0' && raw[i+1] != '1') {
					return nil, fmt.Errorf("invalid pointer token %q", raw)
				}
				i++
			}
		}

		token := strings.ReplaceAll(
			strings.ReplaceAll(raw, "~1", "/"), "~0", "~")

		switch node := value.(type) {
		case map[string]any:
			var ok bool
			value, ok = node[token]
			if !ok {
				return nil, fmt.Errorf("pointer key %q not found", token)
			}

		case []any:
			// JSON Pointer array indexes cannot have leading zeros.
			if token == "" || (len(token) > 1 && token[0] == '0') {
				return nil, fmt.Errorf("invalid array index %q", token)
			}
			for _, ch := range token {
				if ch < '0' || ch > '9' {
					return nil, fmt.Errorf("invalid array index %q", token)
				}
			}
			i, err := strconv.Atoi(token)
			if err != nil || i < 0 || i >= len(node) {
				return nil, fmt.Errorf("array index %q out of bounds", token)
			}
			value = node[i]

		default:
			return nil, fmt.Errorf("cannot traverse pointer token %q", token)
		}
	}
	return value, nil
}

func (r *resolver) expand(
	value any, base *url.URL, active map[string]bool,
) (any, error) {
	switch node := value.(type) {
	case map[string]any:
		if refValue, ok := node["$ref"]; ok {
			ref, ok := refValue.(string)
			if !ok {
				return nil, fmt.Errorf("$ref must be a string")
			}

			hasSiblings := len(node) != 1
			if hasSiblings && r.siblings == siblingError {
				return nil, fmt.Errorf(
					"$ref %q has sibling properties; "+
						"select -ref-siblings=ignore, merge, or allOf "+
						"to choose how to handle them", ref)
			}
			if hasSiblings && r.siblings == siblingIgnore {
				// Ignore sibling properties, as older JSON Reference
				// implementations and OpenAPI 3.0 Reference Objects do.
				return r.expandReference(ref, base, active)
			}

			resolved, err := r.expandReference(ref, base, active)
			if err != nil {
				return nil, err
			}
			if !hasSiblings {
				return resolved, nil
			}

			siblingProperties := make(map[string]any, len(node)-1)
			for key, child := range node {
				if key != "$ref" {
					siblingProperties[key] = child
				}
			}
			expandedSiblings, err := r.expand(
				siblingProperties, base, active)
			if err != nil {
				return nil, fmt.Errorf("expanding $ref siblings: %w", err)
			}

			switch r.siblings {
			case siblingMerge:
				target, ok := resolved.(map[string]any)
				if !ok {
					return nil, fmt.Errorf(
						"cannot merge $ref %q siblings into a non-object target",
						ref)
				}
				siblings, ok := expandedSiblings.(map[string]any)
				if !ok {
					return nil, fmt.Errorf(
						"expanded $ref siblings are not an object")
				}
				merged := make(map[string]any, len(target)+len(siblings))
				for key, value := range target {
					merged[key] = value
				}
				for key, value := range siblings {
					merged[key] = value
				}
				return merged, nil

			case siblingAllOf:
				// Preserve the referenced schema and sibling schema as
				// separate constraints, matching modern JSON Schema behavior.
				return map[string]any{
					"allOf": []any{resolved, expandedSiblings},
				}, nil

			default:
				return nil, fmt.Errorf(
					"unsupported $ref sibling semantics %q", r.siblings)
			}
		}

		result := make(map[string]any, len(node))
		for key, child := range node {
			resolved, err := r.expand(child, base, active)
			if err != nil {
				return nil, fmt.Errorf("property %q: %w", key, err)
			}
			result[key] = resolved
		}
		return result, nil

	case []any:
		result := make([]any, len(node))
		for i, child := range node {
			resolved, err := r.expand(child, base, active)
			if err != nil {
				return nil, fmt.Errorf("array index %d: %w", i, err)
			}
			result[i] = resolved
		}
		return result, nil

	default:
		return value, nil
	}
}

func (r *resolver) expandReference(
	ref string, base *url.URL, active map[string]bool,
) (any, error) {
	relative, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}
	target := base.ResolveReference(relative)

	// Canonicalize the already-decoded fragment representation.
	target.RawFragment = ""
	key := target.String()
	if active[key] {
		return nil, fmt.Errorf(
			"recursive reference cannot be fully inlined: %s", key)
	}

	doc, err := r.load(target)
	if err != nil {
		return nil, err
	}
	selected, err := pointer(doc, target.Fragment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	active[key] = true
	result, err := r.expand(selected, documentURL(target), active)
	delete(active, key)
	return result, err
}

func run() error {
	const siblingHelp = `How to handle properties beside $ref:
  error   Reject $ref objects with siblings (the default).
  ignore  Discard siblings and use only the referenced value, as in older
          JSON Reference implementations and OpenAPI 3.0 Reference Objects.
  merge   Shallow-merge siblings into the referenced object; sibling keys
          replace referenced keys with the same name.
  allOf   Keep the referenced value and expanded siblings as separate entries
          in an allOf array, matching modern JSON Schema sibling semantics.`

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [flags] <input.json-or-URL>\n", os.Args[0])
		fs.PrintDefaults()
	}
	siblings := fs.String(
		"ref-siblings",
		string(siblingError),
		siblingHelp,
	)
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("expected exactly one input")
	}

	semantics := siblingSemantics(*siblings)
	switch semantics {
	case siblingError, siblingIgnore, siblingMerge, siblingAllOf:
	default:
		return fmt.Errorf(
			"invalid -ref-siblings value %q; expected error, ignore, merge, or allOf",
			*siblings)
	}

	u, err := inputURL(fs.Arg(0))
	if err != nil {
		return err
	}
	if u.Fragment != "" {
		return fmt.Errorf("input must identify a document, not a fragment")
	}

	r := &resolver{
		docs:     make(map[string]any),
		client:   &http.Client{Timeout: 30 * time.Second},
		siblings: semantics,
	}
	doc, err := r.load(u)
	if err != nil {
		return err
	}
	resolved, err := r.expand(doc, documentURL(u), make(map[string]bool))
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(resolved)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}


// ``` reasoning
// **Clarifying Go strings.SplitSeq**

// I’m thinking about the Go strings.SplitSeq method introduced in version 1.24—though, hmm, I wonder if that's actually relevant here. I want to give an example and explain its properties, especially whether there are any members besides `$ref`, beyond just Go struct fields. I should include a code patch in my explanation and be careful about handling sibling `$ref` strings and nested references. I need to replace the len check and adjust the return accordingly.
// ```

// A **sibling property** is any property in the same JSON object as `$ref`, other than `$ref` itself. For example:

// ```json
// {
//   "$ref": "common.json#/User",
//   "description": "The account owner"
// }
// ```

// Here, `description` is a sibling. Your script rejects these because different JSON-related specifications give them different meanings.

// There isn’t one universally correct merge rule:

// - In modern **JSON Schema**, `$ref` siblings are evaluated alongside the referenced schema; simply overwriting fields may change their meaning.
// - In some **OpenAPI** contexts, sibling properties are ignored or restricted.
// - For a general JSON document, you need to choose a policy.

// If a **shallow overlay** is appropriate for your data, expand the referenced value and the siblings separately, then let sibling properties override matching keys in the referenced object:

// ```go
// result, err := r.expand(selected, documentURL(target), active)
// delete(active, key)
// if err != nil {
// 	return nil, err
// }

// if len(node) == 1 {
// 	return result, nil
// }

// merged, ok := result.(map[string]any)
// if !ok {
// 	return nil, fmt.Errorf("$ref %q resolves to a non-object with siblings", ref)
// }

// siblings := make(map[string]any, len(node)-1)
// for k, v := range node {
// 	if k != "$ref" {
// 		siblings[k] = v
// 	}
// }

// expanded, err := r.expand(siblings, base, active)
// if err != nil {
// 	return nil, err
// }
// for k, v := range expanded.(map[string]any) {
// 	merged[k] = v
// }
// return merged, nil
// ```

// This is only a generic overlay policy: it **replaces** conflicting top-level properties; it does not recursively merge them. For JSON Schema or OpenAPI, it’s safer to implement the rules for the specific format you’re resolving, or keep the current rejection rather than silently producing incorrect output.
