#!/usr/bin/env python3
"""CLI tool for managing and executing Ollama agent tools."""

import json
import sys
from typing import List, Optional

import typer
from rich.console import Console
from rich.table import Table

from toolbox.registry import tool_registry

app = typer.Typer(help="Ollama Agent Tools CLI")
console = Console()


@app.command()
def list():
    """List all available tools."""
    tools = tool_registry.list_tools()

    if not tools:
        console.print("[yellow]No tools registered[/yellow]")
        return

    table = Table(title="Available Tools")
    table.add_column("Name", style="cyan", no_wrap=True)
    table.add_column("Description", style="green")
    table.add_column("Parameters", style="yellow")

    for tool in tools:
        params = ", ".join(
            f"{name}{'*' if p.required else ''}" for name, p in tool.parameters.items()
        )
        table.add_row(tool.name, tool.description, params or "none")

    console.print(table)


@app.command()
def schema(output: str = typer.Option(None, "--output", "-o", help="Output file path")):
    """Export tool schemas in Ollama format."""
    schemas = tool_registry.to_ollama_format()
    output_json = json.dumps(schemas, indent=2)

    if output:
        with open(output, "w") as f:
            f.write(output_json)
        console.print(f"[green]Schemas exported to {output}[/green]")
    else:
        console.print(output_json)


@app.command()
def execute(
    tool_name: str = typer.Argument(..., help="Name of the tool to execute"),
    params: Optional[List[str]] = typer.Argument(
        default=None, help="Parameters as key=value pairs"
    ),
):
    """Execute a tool with given parameters.

    Example: tool-runner execute calculate expression="2 + 2"
    """
    # Parse parameters
    kwargs = {}
    if params:
        for param in params:
            if "=" not in param:
                console.print(
                    f"[red]Error: Invalid parameter format '{param}'. Use key=value[/red]"
                )
                raise typer.Exit(1)
            key, value = param.split("=", 1)
            kwargs[key] = value

    try:
        result = tool_registry.execute(tool_name, **kwargs)
        console.print("[green]Result:[/green]")
        console.print(result)
    except ValueError as e:
        console.print(f"[red]Error: {e}[/red]")
        raise typer.Exit(1)
    except Exception as e:
        console.print(f"[red]Execution Error: {e}[/red]")
        raise typer.Exit(1)


@app.command()
def info(tool_name: str = typer.Argument(..., help="Name of the tool")):
    """Show detailed information about a tool."""
    tool_data = tool_registry.get_tool(tool_name)
    if not tool_data:
        console.print(f"[red]Tool '{tool_name}' not found[/red]")
        raise typer.Exit(1)

    schema, _ = tool_data

    console.print(f"[bold cyan]{schema.name}[/bold cyan]")
    console.print(f"\n[yellow]Description:[/yellow] {schema.description}")

    if schema.parameters:
        console.print("\n[yellow]Parameters:[/yellow]")
        for name, param in schema.parameters.items():
            required = "[red]*[/red]" if param.required else ""
            console.print(f"  • {name}{required}: {param.description}")
            console.print(f"    Type: {param.type}")
            if param.enum:
                console.print(f"    Options: {', '.join(param.enum)}")
    else:
        console.print("\n[yellow]Parameters:[/yellow] none")


if __name__ == "__main__":
    app()
