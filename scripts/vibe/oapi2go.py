#!/usr/bin/env python3
"""Generate Go structs from an OpenAPI JSON document read from stdin."""

from __future__ import annotations

import argparse
import json
import re
import sys

from collections.abc import Callable, Iterator
from dataclasses import dataclass
from typing import Any

_GO_KEYWORDS = frozenset({
    "break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct",
    "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type",
    "continue", "for", "import", "return",
    "var",
})


def camel_name(value: str) -> str:
    parts = re.findall(r"[A-Za-z0-9]+", value)
    result = "".join(part[0].upper() + part[1:] for part in parts if part)
    if not result:
        return "X"
    if result[0].isdigit():
        result = "X" + result
    return result


#################
# walk function #

@dataclass
class Context:
    key: str
    depth: int


def walk(
    document: Any,
    dispatch: dict[type[Any], Callable[[Context, Any], Any]],
    fallback: Callable[[Context, Any], Any],
    root_key: str = "",
) -> Exception | None:
    def visit(context: Context, value: Any) -> Exception | None:
        result = dispatch.get(type(value), fallback)(context, value)

        if isinstance(result, Exception):
            return result

        if result is None:
            return None

        for child_key, child_value in result:
            child_context = Context(
                key=str(child_key),
                depth=context.depth + 1,
            )
            error = visit(child_context, child_value)
            if error is not None:
                return error

        return None

    return visit(Context(key=root_key, depth=0), document)


###########
# walkers #

def debug(document: Any) -> Exception | None:
    def parse_dict(
        context: Context, value: dict[Any, Any]
    ) -> Iterator[tuple[str, Any]]:
        print(f"{'  ' * context.depth}{{{context.key}}}/{len(value)}")
        return ((str(child_key), child) for child_key, child in value.items())

    def parse_list(
        context: Context, value: list[Any]
    ) -> Iterator[tuple[str, Any]]:
        print(f"{'  ' * context.depth}[{context.key}]/{len(value)}")
        return ((str(index), child) for index, child in enumerate(value))

    def print_literal(context: Context, value: Any) -> None:
        print(f"{'  ' * context.depth} {type(value).__name__} {context.key} -> {value}")

    def parse_failure(context: Context, value: Any) -> Exception:
        print_literal(context, value)
        print(f"unsupported type {type(value).__name__}, aborting")
        print(value)
        return TypeError(f"unsupported type {type(value).__name__}")

    dispatch = {
        dict: parse_dict,
        list: parse_list,
        str: print_literal,
        bool: print_literal,
        int: print_literal,
        type(None): print_literal,
        float: print_literal,
    }

    return walk(document, dispatch, parse_failure, "root")


#################
# main function #

def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--package", required=True, help="Go package name")
    parser.add_argument(
        "--keep",
        help=(
            "regular expression selecting which input JSON paths to keep "
            "(e.g. '^paths\\./models\\.get\\.responses\\.200'); "
            "top-level roots whose JSON path does not match are skipped, "
            "along with all schemas nested inside them"
        ),
    )

    parser.add_argument("--debug", action="store_true", help="enable debug output")
    args = parser.parse_args()

    if not re.fullmatch(r"[A-Za-z_][A-Za-z0-9_]*", args.package):
        parser.error("--package must be a valid Go package identifier")
    if args.package == "_" or args.package in _GO_KEYWORDS:
        parser.error("--package must be a valid Go package identifier")

    try:
        document = json.load(sys.stdin)
        if not isinstance(document, dict):
            raise ValueError("the JSON document must be an object")

        keep = re.compile(args.keep) if args.keep else None
        if args.debug:
            debug(document)
        else:
            print("parsing NIY")

    except (json.JSONDecodeError, ValueError, TypeError, re.error) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

