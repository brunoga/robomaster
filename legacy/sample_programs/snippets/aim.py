# Robomaster S1 camera field of view information.
ROBOT_CAMERA_HORIZONTAL_FOV = 96
ROBOT_CAMERA_VERTICAL_FOV   = 54

# Aim return codes.
AIM_ERROR       = 0
AIM_IN_PROGRESS = 1
AIM_DONE        = 2

# Aims the Robomaster S1 gimbal to the given coordinates.
#
# There are 2 possible operating modes supported:
#
# Direct Mode computes the angle between the source (sight) and destination
# (detected object) positions and uses that angle to directly move the gimbal to
# the target. In this mode only dst_x and dst_y should be provided.
#
# PID Mode computes the delta between the source (sight) and destination
# (detected object) positions and feeds this delta as errors to the given PID
# controllers. The gimbal turn speed is then set based on the PID controllers
# output.
#
# Return values:
#
# AIM_ERROR indicates invalid parameters were provided.
# AIM_IN_PROGRESS indicates we are currently trying to aim at the destination
#   position but did not get a lock yet.
# AIM_DONE indicates we are now locked to the destination position.
#
# Note that no matter which method is used, at least 2 passes are required to
# get an AIM_DONE as the first pass will always return AIM_IN_PROGRESS (unless
# the robot was already pointing directly to the destination position).
#
# Actions that require a target lock should only be done when AIM_DONE is
# returned.
def Aim(dst_x, dst_y, pid_yaw = None, pid_pitch = None):
    if dst_x < 0.0 or dst_x > 1.0 or dst_y < 0.0 or dst_y > 1.0:
        # Invalid dst_x or dst_y.
        return AIM_ERROR

    if ((pid_yaw is not None and pid_pitch is None) or
            (pid_pitch is not None and pid_yaw is None)):
        # Only one of pid_yaw and pid_pitch was provided.
        return AIM_ERROR

    # Obtain sight position. This takes into account sight calibration.
    sight_info = media_ctrl.get_sight_bead_position()
    src_x = sight_info[0]
    src_y = sight_info[1]

    # Compute deltas between source and destination.
    delta_x = dst_x - src_x
    delta_y = src_y - dst_y

    if abs(delta_x) <= 0.1 and abs(err_y) <= 0.1:
        # We are centered in the target already. There is nothing else to do.
        if not None pid_yaw:
            # We are in PID mode. Stop gimbal rotation that might still be in
            # progress.
            gimbal_ctrl.rotate_with_speed(0, 0)
        return AIM_DONE

    if pid_yaw is not None:
        # PID mode.

        # Set error in the PID controllers.
        pid_yaw.set_error(delta_x)
        pid_pitch.set_error(delta_y)

        # Set gimbal rotation speed based on PID controllers output.
        gimbal_ctrl.rotate_with_speed(pid_yaw.get_output(),
                pid_pitch.get_output())
    else:
        # Direct mode.

        # Get current gimbal yaw and pitch angles.
        gimbal_yaw_angle = gimbal_ctrl.get_axis_angle(rm_define.gimbal_axis_yaw)
        gimbal_pitch_angle = gimbal_ctrl.get_axis_angle(
                rm_define.gimbal_axis_pitch)

        # Compute deltas between source and destination angles.
        delta_yaw_angle = ROBOT_CAMERA_HORIZONTAL_FOV * (delta_x)
        delta_pitch_angle = ROBOT_CAMERA_VERTICAL_FOV * (delta_y)

        # Move gimbal so the sight points directly to the target.
        gimbal_ctrl.angle_ctrl(gimbal_yaw_angle + delta_yaw_angle,
                gimbal_pitch_angle + delta_pitch_angle)

    return AIM_IN_PROGRESS
