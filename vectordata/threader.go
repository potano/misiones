// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata
        

type gatheredSegment struct {
	obj mapItemType
	paths []gatheredPath
	lat1, long1, lat2, long2 float64
}

type gatheredPath struct {
	path *map_locationType
	startPoint, endPoint int
}



func (seg *mapSegmentType) threadPaths() (*gatheredSegment, error) {
	allPaths := []mapItemType{}
	for _, item := range seg.paths {
		switch item := item.(type) {
		case *map_locationType:
			allPaths = append(allPaths, item)
		case *map_referenceAggregateType:
			for _, mem := range item.targets {
				_, is := mem.(*map_locationType)
				if !is {
					return nil, item.Error("%s is not a location", item.Name())
				}
				allPaths = append(allPaths, mem)
			}
		default:
			return nil, item.Error("'%s' not allowed here", item.Name())
		}
	}
	paths := make([]gatheredPath, 0, len(allPaths))
	var startLat, startLong, nextLat, nextLong float64
	var pendingPath *pendingPathType
	var started bool
	var item, prevItem mapItemType
	for _, item = range allPaths {
		var isPath bool
		switch item.ItemType() {
		case mitPath:
			isPath = true
		case mitPoint, mitMarker, mitCircle:
		default:
			return nil, item.Error("%s is not legal in a segment",
				typeMapToName[item.ItemType()])
		}
		loc := item.(*map_locationType)
		if isPath {
			if pendingPath == nil {
				pendingPath = newPendingPath(loc)
				if started {
					if !pendingPath.setStartpoint(nextLat, nextLong) {
						goto noConnectError
					}
				}
			} else {
				if !pendingPath.sendEndpointToNextPath(loc) {
					goto noConnectError
				}
				if !started {
					startLat, startLong = pendingPath.getStartpoint()
					started = true
				}
				nextLat, nextLong = pendingPath.getEndpoint()
				paths = append(paths, pendingPath.flush())
				pendingPath = newPendingPath(loc)
			}
		} else {
			waypointLat, waypointLong := loc.location[0], loc.location[1]
			if pendingPath == nil {
				if started {
					if !isSamePoint(waypointLat, waypointLong,
							nextLat, nextLong) {
						goto noConnectError
					}
				} else {
					startLat, startLong = waypointLat, waypointLong
					started = true
				}
			} else {
				if !pendingPath.setEndpoint(waypointLat, waypointLong) {
					goto noConnectError
				}
				if !started {
					startLat, startLong = pendingPath.getStartpoint()
					started = true
				}
				paths = append(paths, pendingPath.flush())
				pendingPath = nil
			}
			nextLat, nextLong = waypointLat, waypointLong
		}
		prevItem = item
	}
	if pendingPath != nil {
		if started {
			if !pendingPath.setStartpoint(nextLat, nextLong) {
				goto noConnectError
			}
		} else {
			startLat, startLong = pendingPath.getStartpoint()
			started = true
		}
		nextLat, nextLong = pendingPath.getEndpoint()
		paths = append(paths, pendingPath.flush())
	}
	if len(paths) == 0 {
		return nil, seg.Error("segment '%s' is empty", seg.Name())
	}
	return &gatheredSegment{seg, paths, startLat, startLong, nextLat, nextLong}, nil

	noConnectError:
	if prevItem != nil {
		return nil, seg.Error("%s '%s' does not connect with %s '%s' in segment '%s'",
			typeMapToName[item.ItemType()], item.Name(),
			typeMapToName[prevItem.ItemType()], prevItem.Name(), seg.Name())
	}
	return nil, seg.Error("%s '%s' does not connect with segment '%s'",
		typeMapToName[item.ItemType()], item.Name(), seg.Name())
}




func (seg *gatheredSegment) reverse() {
	reversedPaths := make([]gatheredPath, len(seg.paths))
	pos := len(seg.paths)
	for _, pth := range seg.paths {
		pth.startPoint, pth.endPoint = pth.endPoint, pth.startPoint
		pos--
		reversedPaths[pos] = pth
	}
	seg.paths = reversedPaths
	seg.lat1, seg.lat2 = seg.lat2, seg.lat1
	seg.long1, seg.long2 = seg.long2, seg.long1
}



func (gp gatheredPath) points() (locationPairs, int, bool) {
	startPoint, endPoint := gp.startPoint, gp.endPoint
	var forward bool
	if endPoint < startPoint {
		startPoint, endPoint = endPoint, startPoint
	} else {
		forward = true
	}
	return gp.path.location[startPoint:endPoint+2], startPoint >> 1, forward
}




type pendingPathType struct {
	path *map_locationType
	startPoint, endPoint int
}

func newPendingPath(path *map_locationType) *pendingPathType {
	return &pendingPathType{path, -1, -1}
}

func (pp *pendingPathType) setStartpoint(lat, long float64) bool {
	loc := pp.path.location
	pp.startPoint = loc.indexOfPoint(lat, long)
	if pp.startPoint == len(loc)-2 && pp.endPoint < 0 {
		//This handles an important special case:  path at the end of the list
		// with points listed in reverse of of previous path.  We can disambiguate this
		// special case, but a waypoint midway in path needs user's explicit indication
		// of a waypoint to end the segment
		pp.endPoint = 0
	}
	return pp.startPoint >= 0
}

func (pp *pendingPathType) setEndpoint(lat, long float64) bool {
	pp.endPoint = pp.path.location.indexOfPoint(lat, long)
	if pp.endPoint == 0 && pp.startPoint < 0 {
		//Corresponding special case for direction-flipping waypoint at end of segment
		pp.startPoint = len(pp.path.location) - 2
	}
	return pp.endPoint >= 0
}

func (pp *pendingPathType) sendEndpointToNextPath(path *map_locationType) bool {
	loc := pp.path.location
	lastIndex := len(loc) - 2
	if path.location.haveMatchingEndpoint(loc[0], loc[1]) {
		pp.endPoint = 0
		if pp.startPoint < 0 {
			pp.startPoint = lastIndex
		}
		return true
	}
	if path.location.haveMatchingEndpoint(loc[lastIndex], loc[lastIndex+1]) {
		pp.endPoint = lastIndex
		if pp.startPoint < 0 {
			pp.startPoint = 0
		}
		return true
	}
	return false
}

func (pp *pendingPathType) getStartpoint() (float64, float64) {
	loc := pp.path.location
	var startPoint int
	if pp.startPoint >= 0 {
		startPoint = pp.startPoint
	}
	return loc[startPoint], loc[startPoint + 1]
}

func (pp *pendingPathType) getEndpoint() (float64, float64) {
	loc := pp.path.location
	var endPoint int
	if pp.endPoint < 0 {
		endPoint = len(loc) - 2
	} else {
		endPoint = pp.endPoint
	}
	return loc[endPoint], loc[endPoint + 1]
}

func (pp *pendingPathType) flush() gatheredPath {
	loc := pp.path.location
	startPoint := pp.startPoint
	if startPoint < 0 {
		startPoint = 0
	}
	endPoint := pp.endPoint
	if endPoint < 0 {
		endPoint = len(loc) - 2
	}
	return gatheredPath{pp.path, startPoint, endPoint}
}




func (route *mapRouteType) threadSegments() ([]*gatheredSegment, error) {
	segments := []*gatheredSegment{}
	for _, seg := range route.segments {
		switch item := seg.(type) {
		case *mapSegmentType:
			gathered, err := item.threadPaths()
			if err != nil {
				return nil, err
			}
			segments = append(segments, gathered)
		case *map_referenceAggregateType:
			for _, mem := range item.targets {
				switch seg := mem.(type) {
				case *mapSegmentType:
					gathered, err := seg.threadPaths()
					if err != nil {
						return nil, err
					}
					segments = append(segments, gathered)
				default:
					return nil, mem.Error("%s not allowed in segments",
						typeMapToName[mem.ItemType()])
				}
			}
		default:
			// Ignore other object types in routes: they don't affect route length
		}
	}
	var nextLat, nextLong float64
	for segX, seg := range segments {
		if segX > 0 {
			var ok, reverse bool
			if isSamePoint(nextLat, nextLong, seg.lat1, seg.long1) {
				ok = true
			} else if isSamePoint(nextLat, nextLong, seg.lat2, seg.long2) {
				ok = true
				reverse = true
			} else if segX == 1 {
				lat0, long0 := segments[0].lat1, segments[0].long1
				if isSamePoint(lat0, long0, seg.lat1, seg.long1) {
					ok = true
				} else if isSamePoint(lat0, long0, seg.lat2, seg.long2) {
					ok = true
					reverse = true
				}
				segments[0].reverse();
			}
			if !ok {
				return nil, route.Error(
					"segment '%s' does not connect with route '%s'",
					seg.obj.Name(), route.Name())
			}
			if reverse {
				seg.reverse()
			}
		}
		nextLat, nextLong = seg.lat2, seg.long2
	}
	return segments, nil
}




func (path *map_locationType) pathAsGatheredSegment() *gatheredSegment {
	startPoint, endPoint := 0, len(path.location) - 2
	loc := path.location
	return &gatheredSegment{nil, []gatheredPath{{path, startPoint, endPoint}},
		loc[0], loc[1], loc[endPoint-2], loc[endPoint-1]}
}

