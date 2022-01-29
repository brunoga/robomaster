def start():
    # Enable manual control of chassis and gimbal.
    chassis_ctrl.enable_stick_overlay()
    gimbal_ctrl.enable_stick_overlay()
   
    # Enable detection of S1 robots.
    vision_ctrl.enable_detection(rm_define.vision_detection_car)

    # Fire one bead per trigger. This does not affect IR firing.
    gun_ctrl.set_fire_count(1)
   
    # Create PID controllers for pitch and yaw.
    pidPitch = rm_ctrl.PIDCtrl()
    pidYaw = rm_ctrl.PIDCtrl()

    # Set contoller parameters.
    pidPitch.set_ctrl_params(90,0,3)
    pidYaw.set_ctrl_params(120,0,5)
   
    # Smaller or equal error values mean we are aiming straight at robot.
    errThreshold = 0.1
   
    while True:
        # Get list of detected S1 robots.
        robotList=RmList(vision_ctrl.get_car_detection_info())
        if robotList[1] > 0:
            # We found at least one robot.
            
            # Set travel mode to free mode so we can automatically rotate the
            # gimbal.
            robot_ctrl.set_mode(rm_define.robot_mode_free)
            
            # Get coordinates to the center of the first detected S1 robot. This
            # will be our target.
            X = robotList[2]
            Y = robotList[3]
            
            # Compute errors in the X and Y axes.
            errX = X - 0.5
            errY = 0.5 - Y
            
            if abs(errX) <= errThreshold and abs(errY) <= errThreshold:
                # We are centered in our target.
               
                # Stop rotating so we do not move past it.
                gimbal_ctrl.rotate_with_speed(0,0)
               
                # Fire!
                gun_ctrl.fire_once()
               
                # Sleep after firing when not tracking an S1 robot.
                time.sleep(0.5)
               
            else:
                # Set errors into our PID controllers.
                pidYaw.set_error(errX)
                pidPitch.set_error(errY)
               
                # Rotate gimbal to the center of the S1 Robot based on the PID
                # controllers.
                gimbal_ctrl.rotate_with_speed(pidYaw.get_output(),
                        pidPitch.get_output())
               
        else:
            # No robot in sight. Stop rotating.
            gimbal_ctrl.rotate_with_speed(0,0)
            robot_ctrl.set_mode(rm_define.robot_mode_chassis_follow)

            # Sleep when not tracking an S1 robot.
            time.sleep(0.1)
