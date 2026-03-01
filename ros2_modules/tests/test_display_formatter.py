import unittest

from ros2_modules.display_controller.status_formatter import format_display_lines


class TestDisplayFormatter(unittest.TestCase):
    def test_formats_default_status(self) -> None:
        lines = format_display_lines({})
        self.assertEqual(lines[0], "ARM STATUS")
        self.assertEqual(lines[1], "STS:---")
        self.assertEqual(lines[2], "B:  0.0")

    def test_formats_with_sequence_state(self) -> None:
        lines = format_display_lines(
            {
                "is_homed": True,
                "is_enabled": True,
                "is_moving": True,
                "joints": {"base": 10.2, "shoulder": 20.4, "elbow": 30.6},
            },
            "accepted",
        )
        self.assertEqual(lines[1], "STS:HEM")
        self.assertEqual(lines[-1], "SEQ:accepted")


if __name__ == "__main__":
    unittest.main()
