from logging import DEBUG, ERROR, INFO, WARNING

from dosh.logger import get_logger, set_verbosity


def test_verbosity_level():
    levels = [
        (0, ERROR),
        (1, WARNING),
        (2, INFO),
        (3, DEBUG),
        (4, DEBUG),
        (5, DEBUG),
    ]
    for verbosity_level, log_level in levels:
        set_verbosity(verbosity_level)
        assert get_logger().level == log_level
