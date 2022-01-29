# Starts, pauses, or stops the timer.
def timer_ctrl(behavior_enum):
    print(f'tools.timer_ctrl({behavior_enum})')

# Obtains the total time elapsed from when the timer first started to the
# current time (in seconds).
def timer_current():
    print('tools.timer_current()')
    return 0.0

# Obtains the program running time (in seconds).
def run_time_of_program():
    print('tools.run_time_of_program()')
    return 0.0

# Acquires current time information including the year, month, day, hour,
# minute, and second.
def get_localtime(time_enum):
    print(f'tools.get_local_time({time_enum})')
    return 0

# Indicates the total time elapsed from when the robot started running up to the
# current time (in seconds).
def get_unixtime():
    print('tools.get_unixtime()')
    return 0.0

