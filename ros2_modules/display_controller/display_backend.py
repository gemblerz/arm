"""Display backend abstraction for SSD1306-compatible OLEDs."""

from __future__ import annotations

from typing import Sequence

try:
    from luma.core.interface.serial import i2c
    from luma.oled.device import ssd1306
    from PIL import Image, ImageDraw

    _LUMA_AVAILABLE = True
except ImportError:  # pragma: no cover - optional hardware dependency
    _LUMA_AVAILABLE = False


class SSD1306Display:
    """Render text lines to an SSD1306 display with a mock fallback."""

    def __init__(self, *, use_mock: bool = False, address: int = 0x3C, port: int = 1) -> None:
        self._mock = use_mock or not _LUMA_AVAILABLE
        self._last_frame: list[str] = []
        self._device = None
        if not self._mock:
            serial = i2c(port=port, address=address)
            self._device = ssd1306(serial)

    @property
    def last_frame(self) -> list[str]:
        """Return the last rendered frame for testing and diagnostics."""
        return self._last_frame.copy()

    def render_lines(self, lines: Sequence[str]) -> None:
        """Render up to 6 lines suitable for a 128x64 OLED."""
        frame = [str(line)[:21] for line in lines[:6]]
        self._last_frame = frame

        if self._mock:
            return

        image = Image.new("1", (128, 64))
        drawer = ImageDraw.Draw(image)
        for index, line in enumerate(frame):
            drawer.text((0, index * 10), line, fill=255)
        self._device.display(image)
