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


#######################
# OpenAPI dataclasses #

@dataclass
class Reference:
    ref: str
    summary: str | None = None
    description: str | None = None


@dataclass
class Contact:
    name: str | None = None
    url: str | None = None
    email: str | None = None


@dataclass
class License:
    name: str
    identifier: str | None = None
    url: str | None = None


@dataclass
class Info:
    title: str
    version: str
    summary: str | None = None
    description: str | None = None
    terms_of_service: str | None = None
    contact: Contact | None = None
    license: License | None = None


@dataclass
class ServerVariable:
    enum: list[str] | None = None
    default: str = ""
    description: str | None = None


@dataclass
class Server:
    url: str
    description: str | None = None
    variables: dict[str, ServerVariable] | None = None


@dataclass
class Tag:
    name: str
    description: str | None = None
    external_docs: dict[str, str] | None = None


@dataclass
class ExternalDocumentation:
    url: str
    description: str | None = None


@dataclass
class Schema:
    ref: str | None = None
    title: str | None = None
    description: str | None = None
    type: str | list[str] | None = None
    format: str | None = None
    properties: dict[str, Schema] | None = None
    items: Schema | None = None
    additional_properties: bool | Schema | None = None
    all_of: list[Schema] | None = None
    one_of: list[Schema] | None = None
    any_of: list[Schema] | None = None
    not_schema: Schema | None = None
    enum: list[Any] | None = None
    required: list[str] | None = None
    nullable: bool | None = None
    default: Any = None
    example: Any = None
    examples: dict[str, Any] | None = None
    minimum: int | float | None = None
    maximum: int | float | None = None
    exclusive_minimum: int | float | bool | None = None
    exclusive_maximum: int | float | bool | None = None
    min_length: int | None = None
    max_length: int | None = None
    pattern: str | None = None
    min_items: int | None = None
    max_items: int | None = None
    unique_items: bool | None = None
    min_properties: int | None = None
    max_properties: int | None = None
    deprecated: bool | None = None
    read_only: bool | None = None
    write_only: bool | None = None
    extensions: dict[str, Any] | None = None


@dataclass
class Example:
    summary: str | None = None
    description: str | None = None
    value: Any = None
    external_value: str | None = None


@dataclass
class Encoding:
    content_type: str | None = None
    headers: dict[str, "Header"] | None = None
    style: str | None = None
    explode: bool | None = None
    allow_reserved: bool | None = None


@dataclass
class MediaType:
    schema: Schema | Reference | None = None
    example: Any = None
    examples: dict[str, Example | Reference] | None = None
    encoding: dict[str, Encoding] | None = None


@dataclass
class Header:
    description: str | None = None
    required: bool | None = None
    deprecated: bool | None = None
    style: str | None = None
    explode: bool | None = None
    schema: Schema | Reference | None = None
    example: Any = None
    examples: dict[str, Example | Reference] | None = None
    content: dict[str, MediaType] | None = None


@dataclass
class Parameter:
    name: str
    location: str
    description: str | None = None
    required: bool | None = None
    deprecated: bool | None = None
    allow_empty_value: bool | None = None
    style: str | None = None
    explode: bool | None = None
    allow_reserved: bool | None = None
    schema: Schema | Reference | None = None
    example: Any = None
    examples: dict[str, Example | Reference] | None = None
    content: dict[str, MediaType] | None = None


@dataclass
class RequestBody:
    description: str | None = None
    content: dict[str, MediaType] | None = None
    required: bool | None = None


@dataclass
class Response:
    description: str
    headers: dict[str, Header | Reference] | None = None
    content: dict[str, MediaType] | None = None
    links: dict[str, Any] | None = None


@dataclass
class SecurityRequirement:
    schemes: dict[str, list[str]]


@dataclass
class Operation:
    responses: dict[str, Response | Reference]
    tags: list[str] | None = None
    summary: str | None = None
    description: str | None = None
    external_docs: ExternalDocumentation | None = None
    operation_id: str | None = None
    parameters: list[Parameter | Reference] | None = None
    request_body: RequestBody | Reference | None = None
    callbacks: dict[str, Any] | None = None
    deprecated: bool | None = None
    security: list[SecurityRequirement] | None = None
    servers: list[Server] | None = None
    extensions: dict[str, Any] | None = None


@dataclass
class PathItem:
    ref: str | None = None
    summary: str | None = None
    description: str | None = None
    get: Operation | None = None
    put: Operation | None = None
    post: Operation | None = None
    delete: Operation | None = None
    options: Operation | None = None
    head: Operation | None = None
    patch: Operation | None = None
    trace: Operation | None = None
    servers: list[Server] | None = None
    parameters: list[Parameter | Reference] | None = None
    extensions: dict[str, Any] | None = None


@dataclass
class Components:
    schemas: dict[str, Schema | Reference] | None = None
    responses: dict[str, Response | Reference] | None = None
    parameters: dict[str, Parameter | Reference] | None = None
    examples: dict[str, Example | Reference] | None = None
    request_bodies: dict[str, RequestBody | Reference] | None = None
    headers: dict[str, Header | Reference] | None = None
    security_schemes: dict[str, Any] | None = None
    links: dict[str, Any] | None = None
    callbacks: dict[str, Any] | None = None


@dataclass
class OpenAPIDocument:
    openapi: str
    info: Info
    paths: dict[str, PathItem]
    json_schema_dialect: str | None = None
    servers: list[Server] | None = None
    components: Components | None = None
    security: list[SecurityRequirement] | None = None
    tags: list[Tag] | None = None
    external_docs: ExternalDocumentation | None = None


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

