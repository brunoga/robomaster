def start():
    # Enable manual control of chassis and gimbal.
    chassis_ctrl.enable_stick_overlay()
    gimbal_ctrl.enable_stick_overlay()
    
    # Enable detection of S1 robots.
    vision_ctrl.enable_detection(rm_define.vision_detection_car)

    # Fire one bead per trigger. This does not affect IR firing. 
    gun_ctrl.set_fire_count(1)
    
    # Set gimbal rotation speed to the maximum possible, for fast object
    # tracking.
    gimbal_ctrl.set_rotate_speed(540)
   
    # Set travel mode to free mode so we can automatically rotate the gimbal.
    robot_ctrl.set_mode(rm_define.robot_mode_free)

    # Keep track of the current and previous x and y possitions of the
    # detected robot.
    prevX = 0.0
    prevY = 0.0
    x = 0.0
    y = 0.0

    while True:
        # Get list of detected S1 robots.
        robotList=RmList(vision_ctrl.get_car_detection_info())
        if robotList[1] > 0:
            # We found at least one robot.
           
            # Cache previous x and y values.
            prevX = x
            prevY = y
            
            # Get coordinates to the center of the first
            # detected S1 robot. This will be our target.
            x = robotList[2]
            y = robotList[3]
           
            # Robot detection is currently too slow. Make sure that we do not
            # overshoot the target because we think we did not move. Note there
            # are cases this will fail miserably, but they are unlikelly.
            if abs(x - prevX) < 0.01 and abs(y - prevY) < 0.01:
                # It does not look like we got new values. Get robot
                # info again.
                continue
            
            # Get current sight coordinates.
            sightInfo = media_ctrl.get_sight_bead_position()
            sightX = sightInfo[0]
            sightY = sightInfo[1]
            
            # Obtain Compute gimbal angle in relation to the chassis.
            yawAngle = gimbal_ctrl.get_axis_angle(rm_define.gimbal_axis_yaw)
            pitchAngle = gimbal_ctrl.get_axis_angle(rm_define.gimbal_axis_pitch)
            
            # Compute yaw and pitch angle offsets (i.e. how much we need to
            # move).
            yaw = 96 * (x - sightX)
            pitch = 54 * (sightY - y)
            
            if abs(yaw) > 2 or abs(pitch) > 2:
                # Point gimbal to target.
                gimbal_ctrl.angle_ctrl(yawAngle + yaw, pitchAngle + pitch)
           
                # Fire!
                gun_ctrl.fire_once()

