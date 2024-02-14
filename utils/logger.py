from typing import TextIO, Callable
from enum import Enum
from datetime import datetime


class LogLevel(Enum):
    DEBUG = "DEBUG"
    INFO = "INFO"
    WARNING = "WARNING"
    ERROR = "ERROR"
    FATAL = "FATAL"
    NONE = "NONE"


class OutputType(Enum):
    STD = "STD"
    FILE = "FILE"
    STD_AND_FILE = "STD_AND_FILE"


class Timestamp(Enum):
    DATE = "%d-%m-%Y"
    TIME = "%H:%M:%S"
    DATE_AND_TIME = f"{DATE}_{TIME}"
    NONE = "NONE"


class logger:
    STD_LOG_LEVEL: LogLevel = LogLevel.WARNING
    FILE_LOG_LEVEL: LogLevel = LogLevel.DEBUG
    OUTPUT_TYPE: OutputType = OutputType.STD
    OUTPUT_STREAM: TextIO | None = None
    TIMESTAMP: Timestamp = Timestamp.NONE
    FILE_TIMESTAMP: Timestamp = Timestamp.DATE_AND_TIME

    @classmethod
    def setup_logs(cls, args) -> TextIO | None:
        verbose: bool = args.verbose
        debug: bool = args.debug
        silent: bool = args.silent
        force_silent: bool = args.Silent
        logs: bool = args.logs

        file_output: bool = False
        logs_file = None

        if logs:
            cls.FILE_TIMESTAMP = Timestamp.DATE_AND_TIME
            cls.FILE_LOG_LEVEL = LogLevel.DEBUG
            logs_fp = f"logs_oait_{logger._format_date_time(Timestamp.DATE_AND_TIME)}.txt"
            logs_file = open(logs_fp, 'w')
            cls.OUTPUT_STREAM = logs_file
            file_output = True

        if verbose:
            if debug or silent or force_silent:
                print("FATAL: cannot have more than one log level enable.")
                exit(1)

            cls.STD_LOG_LEVEL = LogLevel.INFO
        if debug:
            if verbose or silent or force_silent:
                print("FATAL: cannot have more than one log level enable.")
                exit(1)

            cls.TIMESTAMP = Timestamp.TIME
            cls.STD_LOG_LEVEL = LogLevel.DEBUG

        if silent:
            if verbose or debug or force_silent:
                print("FATAL: cannot have more than one log level enable.")
                exit(1)

            cls.STD_LOG_LEVEL = LogLevel.FATAL

        if force_silent:
            if verbose or debug or silent:
                print("FATAL: cannot have more than one log level enable.")
                exit(1)

            cls.STD_LOG_LEVEL = LogLevel.NONE

        if file_output:
            cls.OUTPUT_TYPE = OutputType.STD_AND_FILE
        else:
            cls.OUTPUT_TYPE = OutputType.STD

        return logs_file

    @classmethod
    def debug(cls, message: str, method: Callable | None = None) -> None:
        cls._output_log(LogLevel.DEBUG, message, method)

    @classmethod
    def info(cls, message: str, method: Callable | None = None) -> None:
        cls._output_log(LogLevel.INFO, message, method)

    @classmethod
    def warning(cls, message: str, method: Callable | None = None) -> None:
        cls._output_log(LogLevel.WARNING, message, method)

    @classmethod
    def error(cls, message: str, method: Callable | None = None) -> None:
        cls._output_log(LogLevel.ERROR, message, method)

    @classmethod
    def fatal(cls, message: str, method: Callable | None = None) -> None:
        cls._output_log(LogLevel.FATAL, message, method)

    @classmethod
    def _output_log(cls, log_type: LogLevel, message: str, method: Callable | None) -> None:

        if (cls.OUTPUT_TYPE == OutputType.STD or cls.OUTPUT_TYPE == OutputType.STD_AND_FILE) and logger._log_level_applies(cls.STD_LOG_LEVEL, log_type):
            output_fmt: str = logger._format_output(cls.TIMESTAMP, log_type, message, method)
            logger._output_std(output_fmt)

        if (cls.OUTPUT_TYPE == OutputType.FILE or cls.OUTPUT_TYPE == OutputType.STD_AND_FILE) and logger._log_level_applies(cls.FILE_LOG_LEVEL, log_type):
            output_fmt: str = logger._format_output(cls.FILE_TIMESTAMP, log_type, message, method)
            logger._output_file(cls.OUTPUT_STREAM, output_fmt)

    @staticmethod
    def _output_std(output: str) -> None:
        print(output)

    @staticmethod
    def _output_file(output_stream: TextIO | None, output: str) -> None:
        if output_stream:
            output_stream.write(output + '\n')
        else:
            print(f"LOGGER FATAL ERROR: Failed to write to OUTPUT_STREAM ('{output_stream}').")

    @staticmethod
    def _log_level_applies(log_level: LogLevel, log_type: LogLevel) -> bool:
        if log_level == LogLevel.NONE:
            return False

        if log_level == LogLevel.FATAL:
            return log_type == LogLevel.FATAL

        if log_level == LogLevel.ERROR:
            return log_type == LogLevel.FATAL or log_type == LogLevel.ERROR

        if log_level == LogLevel.WARNING:
            return log_type == LogLevel.FATAL or log_type == LogLevel.ERROR or log_type == LogLevel.WARNING

        if log_level == LogLevel.INFO:
            return log_type == LogLevel.FATAL or log_type == LogLevel.ERROR or log_type == LogLevel.WARNING or log_type == LogLevel.INFO

        if log_level == LogLevel.DEBUG:
            return log_type == LogLevel.FATAL or log_type == LogLevel.ERROR or log_type == LogLevel.WARNING or log_type == LogLevel.INFO or log_type == LogLevel.DEBUG

        return False

    @staticmethod
    def _format_output(timestamp: Timestamp, log_type: LogLevel, message: str, method: Callable | None) -> str:

        date_time_fmt = logger._format_date_time(timestamp)
        method_fmt = logger._format_method(method)
        log_type_fmt = logger._format_log_type(log_type)

        timestamp_pad = 0 if timestamp == Timestamp.NONE else 25

        output_fmt = f"{date_time_fmt: <{timestamp_pad}}{method_fmt: <30}{log_type_fmt: <10}{message}"
        return output_fmt 

    @staticmethod
    def _format_log_type(log_type: LogLevel) -> str:
        log_type_fmt: str = f"{log_type.value}:"

        return log_type_fmt

    @staticmethod
    def _format_date_time(timestamp: Timestamp) -> str:
        if timestamp == Timestamp.NONE:
            return ""

        date_time_fmt: str = f"{(datetime.now().strftime(timestamp.value))}"

        return date_time_fmt

    @staticmethod
    def _format_method(method: Callable | None) -> str:
        if method is None:
            return ""

        method_fmt: str = f"({method.__name__})"

        return method_fmt 


