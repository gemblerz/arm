"""Sequence command parsing for the ROS2 executor."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Optional

from ros2_modules.common.arm_protocol import CartesianTarget


@dataclass(frozen=True)
class ParsedCommand:
    """Parsed sequence command."""

    action: str
    target: Optional[CartesianTarget] = None


def parse_command(payload: dict) -> ParsedCommand:
    """Parse a sequence command payload."""
    command = str(payload.get("command", "")).strip().lower()
    if command == "home":
        return ParsedCommand(action="home")
    if command in {"move_cartesian", "move_cartician"}:
        target = payload.get("target", {})
        return ParsedCommand(
            action="move_cartesian",
            target=CartesianTarget(
                x=float(target["x"]),
                y=float(target["y"]),
                z=float(target["z"]),
            ),
        )
    raise ValueError(f"unsupported sequence command: {command}")
