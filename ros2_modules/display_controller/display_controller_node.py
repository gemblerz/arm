"""ROS2 node that renders status topics to SSD1306 display."""

from __future__ import annotations

import json

import rclpy
from rclpy.node import Node
from std_msgs.msg import String

from ros2_modules.display_controller.display_backend import SSD1306Display
from ros2_modules.display_controller.status_formatter import format_display_lines


class DisplayControllerNode(Node):
    """Subscribe to arm status topics and render compact OLED output."""

    def __init__(self) -> None:
        super().__init__("display_controller")
        self._display = SSD1306Display(use_mock=True)
        self._latest_status: dict[str, object] = {}
        self._sequence_state = ""
        self.create_subscription(String, "/arm/status", self._on_status, 10)
        self.create_subscription(String, "/arm/sequence_status", self._on_sequence_status, 10)
        self.get_logger().info("Display controller started")

    def _on_status(self, msg: String) -> None:
        self._latest_status = json.loads(msg.data)
        self._refresh()

    def _on_sequence_status(self, msg: String) -> None:
        payload = json.loads(msg.data)
        self._sequence_state = str(payload.get("state", ""))
        self._refresh()

    def _refresh(self) -> None:
        lines = format_display_lines(self._latest_status, self._sequence_state)
        self._display.render_lines(lines)


def main() -> None:
    rclpy.init()
    node = DisplayControllerNode()
    try:
        rclpy.spin(node)
    finally:
        node.destroy_node()
        rclpy.shutdown()


if __name__ == "__main__":
    main()
