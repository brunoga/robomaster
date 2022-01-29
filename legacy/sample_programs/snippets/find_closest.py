# Return the bounding box information (X, Y, W, H) for the closest detected
# robot.
def FindClosestRobot(robotDetectionInfo):
    return FindClosest(robotDetectionInfo, 4)

# Return the bounding box information (ID, X, Y, W, H) for the closest detected
# vision marker.
def FindClosestVisionMarker(markerDetectionInfo)
    return FindClosest(robotDetectionInfo, 5)

# Simple algorithm to find the closest detected object. It simply iterates
# through all of the detected objects and returns the one with the bigger
# height.
#
# Other than the detection info, it takes as parameter the expected number of
# items per detected object as it varies for a detected robot or a detected
# vision marker, for example.
def FindClosest(detectionInfo, numEntriesPerObject):
    numEntries = len(robotDetectionInfo) - 1  # Ignore size entry.
    numObjects = numEntries // numEntriesPerObject
    if numObjects != detectionInfo[0]:
        # Got an unexpected number of entries.
        return None

    modulo = numEntries % numEntriesPerObject
    if modulo != 0:
        # Got incomplete number of entries.
        return None

    closestHeight = 0.0  # Impossible height.
    closestIndex = 1 # Defaults to first Robot detected.

    # Check height of the bounding box of each detected robot. Return the one
    # with the biggest height.
    for i in range(1, len(detectionInfo) - 1, numEntriesPerObject):
        objectHeight = detectionInfo[i + numEntriesPerObject - 1]
        if objectHeight > closestHeight:
            # Found a bigger height.
            closestHeight = objectHeight
            closestIndex = i

    # Return only the relevant info about the selected object.
    return detectionInfo[closestIndex:closestIndex + numEntriesPerObject - 1]
