import unittest

from ros2_modules.sequence_executor.sequence_parser import parse_command


class TestSequenceParser(unittest.TestCase):
    def test_parse_home(self) -> None:
        cmd = parse_command({"command": "home"})
        self.assertEqual(cmd.action, "home")
        self.assertIsNone(cmd.target)

    def test_parse_cartesian_with_legacy_typo_alias(self) -> None:
        cmd = parse_command(
            {"command": "move_cartician", "target": {"x": 0.1, "y": 0.2, "z": 0.3}}
        )
        self.assertEqual(cmd.action, "move_cartesian")
        self.assertAlmostEqual(cmd.target.x, 0.1)
        self.assertAlmostEqual(cmd.target.y, 0.2)
        self.assertAlmostEqual(cmd.target.z, 0.3)

    def test_parse_unknown_raises(self) -> None:
        with self.assertRaises(ValueError):
            parse_command({"command": "dance"})


if __name__ == "__main__":
    unittest.main()
