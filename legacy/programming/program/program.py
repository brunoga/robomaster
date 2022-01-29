import media_ctrl.media_ctrl as media_ctrl
import rm_define.rm_define as rm_define
import tools.tools as tools

# Sample usage of the Robomaster S1 stubs. start() is called automatically. Note
# the program below should be a valid Robomaster S1 program (even if it might
# not do anything useful).
def start():
    rm_define.robot_set_mode(rm_define.robot_mode_free)

    tools.timer_ctrl(rm_define.timer_start)

    elapsed = tools.timer_current()
    print(elapsed)

    tools.timer_ctrl(rm_define.timer_stop)

    media_ctrl.zoom_value_update(2)

