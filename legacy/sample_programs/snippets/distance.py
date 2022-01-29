# Object distance via triangle similarity.
#
# Usage:
#
# There are 3 pieces of information we need to figure out the distance to a
# knonw object in an image:
#
# Known height or width: This is the actual height or width of the object we
#   want to measure the distance to. For our purposes (a land robot) it will
#   almost always be better to use height instead of width as the aparent width 
#   of an object will change depending on its orientation in relationship to the
#   camera (for example, a front facing Robomaster S1 will have a know width of 
#   240 mm while the same robot at the same distance but with the side facing
#   the camera will have a known width of 320 mm).
# Height or width of the object bounding box: This is automatically provided by
#   the Robomaster API as the H or W attribute of an identified object. This
#   will always be a number in the interval [0, 1]. The dimension used here
#   (height or width) has to be the same used above.
# Camera focal length: This can be computed based on the information above plus
#   the actual known distance that resulted in the bounding box dimensions used.
#
# The first thing we need to do is compute the focal length. To do that you need
# to put an object with a known dimension (again, known height is better) at a
# known distance of the camera. Make sure the Robomaster S1 is pointing straight
# at it and check the H or W values. then simply do:
#
#   focalLength = FocalLength(knownHeight, knownDistance, height)
#
# You only need to do this once. After you figured out the value for your camera
# you can simply use that value directly. Fortunately the Robomaster S1 camera
# focal length is even easier to compute due to the coordinate system used (see
# related constant below).
#
# Then, to actually compute the distance, just do:
#
#   distance = Distance(knownHeight, focalLength, height)
#
# The returned distance will be in the same unit as the provided knownHeight.


# Known height of a Robomaster S1 in millimeters and inches.
ROBOT_KNOWN_HEIGHT_MM = 270.0
ROBOT_KNOWN_HEIGHT_IN = 10.6

# Due to the corrdinate system used, the focal length can be inferred directly
# (you can still compute it yourself using FocalLength() to see it matches).
ROBOT_CAMERA_FOCAL_LENGTH = 1.0

# Computes the camera focal length. Here for reference as for the Robomaster S1
# one can simply use the constant above.
def FocalLength(knownHeightOrWidth, knownDistance, heightOrWidth):
    return (heightOrWidth * knownDistance) / knownHeightOrWidth

# Returns the distance to the object in the same unit as the given
# knownHeightOrWidth.
def Distance(knownHeightOrWidth, focalLength, heightOrWidth):
    return (knownHeightOrWidth * focalLength) / heightOrWidth

# Compute distance in millimiters to a detected Robomaster S1 given its
# bounding box height.
def DistanceToRobotMM(height):
    return Distance(ROBOT_KNOWN_HEIGHT_MM, ROBOT_CAMERA_FOCAL_LENGTH,
            height)

# Compute distance in inches to a detected Robomaster S1 given its
# bounding box height).
def DistanceToRobotIn(height):
    return Distance(ROBOT_KNOWN_HEIGHT_IN, ROBOT_CAMERA_FOCAL_LENGTH,
            height)

