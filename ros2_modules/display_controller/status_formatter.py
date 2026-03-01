"""Formatting helpers for display text payloads."""

from __future__ import annotations

from typing import Any


def _to_bool_flag(value: Any, enabled: str, disabled: str) -> str:
    return enabled if bool(value) else disabled


def format_display_lines(status: dict[str, Any], sequence_state: str = "") -> list[str]:
    """Format arm state payload into compact OLED display lines."""
    joints = status.get("joints") or {}
    base = float(joints.get("base", 0.0))
    shoulder = float(joints.get("shoulder", 0.0))
    elbow = float(joints.get("elbow", 0.0))

    flags = "".join(
        [
            _to_bool_flag(status.get("is_homed"), "H", "-"),
            _to_bool_flag(status.get("is_enabled", True), "E", "-"),
            _to_bool_flag(status.get("is_moving"), "M", "-"),
        ]
    )

    lines = [
        "ARM STATUS",
        f"STS:{flags}",
        f"B:{base:5.1f}",
        f"S:{shoulder:5.1f}",
        f"E:{elbow:5.1f}",
    ]
    if sequence_state:
        lines.append(f"SEQ:{sequence_state}")
    return lines
