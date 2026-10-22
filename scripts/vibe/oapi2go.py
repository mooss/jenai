#!/usr/bin/env python3
"""Generate Go structs from an OpenAPI JSON document read from stdin."""

from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass, field
from typing import Any


@dataclass
class IRType:
    """Language-neutral description of a type."""

    kind: str
    nullable: bool = False
    name: str = ""
    item: IRType | None = None
    map_value: IRType | None = None


@dataclass
class IRField:
    name: str
    description: str
    optional: bool
    type: IRType


@dataclass
class IRStruct:
    name: str
    description: str
    fields: list[IRField] = field(default_factory=list)


@dataclass
class IRAlias:
    name: str
    description: str
    type: IRType


@dataclass
class IR:
    structs: list[IRStruct] = field(default_factory=list)
    aliases: list[IRAlias] = field(default_factory=list)


_GO_KEYWORDS = frozenset(
    {
        "break",
        "default",
        "func",
        "interface",
        "select",
        "case",
        "defer",
        "go",
        "map",
        "struct",
        "chan",
        "else",
        "goto",
        "package",
        "switch",
        "const",
        "fallthrough",
        "if",
        "range",
        "type",
        "continue",
        "for",
        "import",
        "return",
        "var",
    }
)


def camel_name(value: str) -> str:
    parts = re.findall(r"[A-Za-z0-9]+", value)
    result = "".join(part[0].upper() + part[1:] for part in parts if part)
    if not result:
        return "X"
    if result[0].isdigit():
        result = "X" + result
    return result

def descriptions(schema: dict[str, Any]) -> str:
    lines = [str(schema["description"])] if schema.get("description") else []
    for keyword in ("allOf", "oneOf", "anyOf"):
        if keyword in schema:
            # This might not always be this exact layout, sorry if it crashes:
            parts = schema[keyword]
            lines.append(parts[1]['description']) # Textual description.
            lines.append(f"{keyword} composition:")
            lines.extend(json.dumps(parts[0], indent=2).splitlines()) # allOf/oneOf/anyOf spec.
    return "\n".join(lines)


def reject_refs(value: Any, path: str = "$") -> None:
    if isinstance(value, dict):
        if "$ref" in value:
            raise ValueError(f"$ref is not supported at {path}")
        for key, child in value.items():
            reject_refs(child, f"{path}.{key}")
    elif isinstance(value, list):
        for index, child in enumerate(value):
            reject_refs(child, f"{path}[{index}]")


def schema_types(schema: dict[str, Any]) -> list[str]:
    value = schema.get("type", [])
    return value if isinstance(value, list) else [value]


def compile_ir(document: dict[str, Any]) -> IR:
    ir = IR()
    alias_sources: dict[str, str] = {}

    def compile_schema(schema: dict[str, Any], hint: str) -> IRType:
        types = schema_types(schema)
        nullable = bool(schema.get("nullable")) or "null" in types
        non_null_types = [value for value in types if value != "null"]

        if "properties" in schema or "object" in non_null_types or "additionalProperties" in schema:
            properties = schema.get("properties", {})
            if not properties:
                if "additionalProperties" not in schema:
                    return IRType("map", nullable, map_value=IRType("any"))
                additional = schema["additionalProperties"]
                if additional is False:
                    pass
                elif isinstance(additional, dict):
                    value_type = compile_schema(additional, hint + "Value")
                    return IRType("map", nullable, map_value=value_type)
                else:
                    return IRType("map", nullable, map_value=IRType("any"))
            if properties and "additionalProperties" in schema and isinstance(
                schema["additionalProperties"], dict
            ):
                raise ValueError(
                    f"additionalProperties with properties is not supported at {hint}"
                )

            struct = IRStruct(camel_name(hint), descriptions(schema))
            ir.structs.append(struct)
            required = set(schema.get("required", []))
            for field_name, field_schema in properties.items():
                field_schema = field_schema or {}
                child_hint = struct.name + camel_name(field_name)
                child_type = compile_schema(field_schema, child_hint)
                struct.fields.append(
                    IRField(
                        name=field_name,
                        description=descriptions(field_schema),
                        optional=field_name not in required,
                        type=child_type,
                    )
                )
            return IRType("object", nullable, name=struct.name)

        if "array" in non_null_types or "items" in schema:
            item_schema = schema.get("items", {})
            return IRType(
                "array",
                nullable,
                item=compile_schema(item_schema, hint + "Item"),
            )

        if "number" in non_null_types:
            return IRType("float64", nullable)
        if "integer" in non_null_types:
            return IRType("int64", nullable)
        if "boolean" in non_null_types:
            return IRType("bool", nullable)
        if "string" in non_null_types:
            return IRType("string", nullable)

        # Composition-only schemas may not declare a type themselves. Use the
        # first variant's type as a best-effort approximation.
        for keyword in ("allOf", "oneOf", "anyOf"):
            variants = schema.get(keyword, [])
            if variants and isinstance(variants[0], dict):
                inferred = compile_schema(variants[0], hint)
                inferred.nullable = inferred.nullable or nullable
                return inferred
        return IRType("any", nullable)

    def add_root(schema: dict[str, Any], name: str, fallback_doc: str = "") -> None:
        result = compile_schema(schema, name)
        if result.kind != "object":
            alias_name = camel_name(name)
            alias_doc = descriptions(schema)
            if fallback_doc:
                schema_doc = alias_doc.strip()
                fallback = fallback_doc.strip()
                already_present = schema_doc == fallback or schema_doc.startswith(
                    fallback + "\n"
                )
                if fallback and not already_present:
                    alias_doc = "\n".join(
                        part for part in (schema_doc, fallback) if part
                    )
            if alias_name in alias_sources:
                raise ValueError(
                    f"alias naming collision: "
                    f"{alias_sources[alias_name]!r} and {name!r} become {alias_name}"
                )
            alias_sources[alias_name] = name
            ir.aliases.append(IRAlias(alias_name, alias_doc, result))
            return
        # Works but generate duplicates.
        # if result.kind == "object" and fallback_doc:
        #     struct = next(s for s in ir.structs if s.name == result.name)
        #     if fallback_doc:
        #         struct.description = "\n".join(
        #             part for part in (struct.description, fallback_doc) if part
        #         )
        if result.kind == "object" and fallback_doc:
            struct = next(s for s in ir.structs if s.name == result.name)

            schema_doc = struct.description.strip()
            fallback = fallback_doc.strip()

            # The schema description may be followed by generated composition
            # comments, so check for an exact description at the start.
            already_present = (
                schema_doc == fallback
                or schema_doc.startswith(fallback + "\n")
            )

            if fallback and not already_present:
                struct.description = "\n".join(
                    part for part in (schema_doc, fallback) if part
                )

    operations: list[tuple[str, dict[str, Any]]] = []
    methods = {"get", "put", "post", "delete", "options", "head", "patch", "trace"}

    if isinstance(document.get("paths"), dict):
        for path, path_item in document["paths"].items():
            if not isinstance(path_item, dict):
                continue
            for method, operation in path_item.items():
                if method.lower() in methods and isinstance(operation, dict):
                    operation_name = operation.get("operationId") or method + path
                    operations.append((camel_name(operation_name), operation))
    else:
        for method, operation in document.items():
            if method.lower() in methods and isinstance(operation, dict):
                operation_name = operation.get("operationId") or method
                operations.append((camel_name(operation_name), operation))

    for operation_name, operation in operations:
        for status, response in operation.get("responses", {}).items():
            if not isinstance(response, dict):
                continue
            content = response.get("content", {})
            media = content.get("application/json")
            if media is None and content:
                media = next(iter(content.values()))
            schema = media.get("schema") if isinstance(media, dict) else None
            if isinstance(schema, dict):
                status_name = "X" + status if status.isdigit() else camel_name(status)
                add_root(
                    schema,
                    operation_name + "Response" + status_name,
                    str(response.get("description", "")),
                )

    components = document.get("components", {}).get("schemas", {})
    if isinstance(components, dict):
        for name, schema in components.items():
            if isinstance(schema, dict):
                add_root(schema, name)

    return ir


def go_base_type(node: IRType) -> str:
    if node.kind == "object":
        return node.name
    if node.kind == "array":
        return "[]" + go_type(node.item or IRType("any"))
    if node.kind == "map":
        return "map[string]" + go_type(node.map_value or IRType("any"))
    return node.kind


def go_type(node: IRType) -> str:
    result = go_base_type(node)
    return "*" + result if node.nullable else result


def comment_lines(text: str, indent: str = "") -> list[str]:
    return [indent + "// " + line if line else indent + "//" for line in text.splitlines()]


def render(ir: IR, package: str, argv: list[str]) -> str:
    seen_structs: dict[str, str] = {}
    for struct in ir.structs:
        if struct.name in seen_structs:
            raise ValueError(
                f"struct naming collision: {struct.name!r} "
                f"({seen_structs[struct.name]} and another declaration)"
            )
        seen_structs[struct.name] = struct.name

        seen_fields: dict[str, str] = {}
        for field_ir in struct.fields:
            go_name = camel_name(field_ir.name)
            if go_name in seen_fields:
                raise ValueError(
                    f"field naming collision in {struct.name}: "
                    f"{seen_fields[go_name]!r} and {field_ir.name!r} become {go_name}"
                )
            seen_fields[go_name] = field_ir.name

    for alias in ir.aliases:
        if alias.name in seen_structs:
            raise ValueError(
                f"struct naming collision: {alias.name!r} "
                f"({seen_structs[alias.name]} and another declaration)"
            )
        seen_structs[alias.name] = alias.name

    output = [
        "// Generated with arguments: "
        + json.dumps(argv, ensure_ascii=False),
        "//revive:disable:var-naming,line-length-limit",
        f"package {package}",
    ]

    for struct in ir.structs:
        output.append("")
        if struct.description:
            output.extend(comment_lines(struct.description))
        output.append(f"type {struct.name} struct {{")
        for field_ir in struct.fields:
            if field_ir.description:
                output.extend(comment_lines(field_ir.description, "\t"))
            field_type = go_type(field_ir.type)
            if field_ir.optional and not field_ir.type.nullable:
                field_type = "*" + field_type
            tag_value = json.dumps(field_ir.name, ensure_ascii=False)
            output.append(
                f"\t{camel_name(field_ir.name)} {field_type} "
                f"`json:{tag_value} yaml:{tag_value}`"
            )
        output.append("}")

    for alias in ir.aliases:
        output.append("")
        if alias.description:
            output.extend(comment_lines(alias.description))
        output.append(f"type {alias.name} {go_type(alias.type)}")

    return "\n".join(output) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--package", required=True, help="Go package name")
    args = parser.parse_args()

    if not re.fullmatch(r"[A-Za-z_][A-Za-z0-9_]*", args.package):
        parser.error("--package must be a valid Go package identifier")
    if args.package == "_" or args.package in _GO_KEYWORDS:
        parser.error("--package must be a valid Go package identifier")

    try:
        document = json.load(sys.stdin)
        if not isinstance(document, dict):
            raise ValueError("the JSON document must be an object")
        reject_refs(document)
        ir = compile_ir(document)
        sys.stdout.write(render(ir, args.package, sys.argv))
    except (json.JSONDecodeError, ValueError, TypeError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
