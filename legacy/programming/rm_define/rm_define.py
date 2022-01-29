# Robot mode "enum".
robot_mode_gimbal_follow = 1
robot_mode_chassis_follow = 2
robot_mode_free = 3

# Timer behavior "enum".
timer_start = 1
timer_stop = 2
timer_reset = 3

# Time "enum".
localtime_year = 1
localtime_month = 2
localtime_day = 3
localtime_hour = 4
localtime_minute = 5
localtime_second = 6

# Sets the travel mode.
#
# * Chassis Lead Mode: The gimbal follows the chassis to rotate along the yaw
#   axis.
# * Gimbal Lead Mode: The chassis follows the gimbal to rotate along the yaw
#   axis.
# * Free Mode: The gimbal and the chassis move without affecting each other.
def robot_set_mode(mode_enum):
    print(f'rm_define.robot_set_mode({mode_enum})')
