// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"

	"potano.misiones/great"
)


func (vd *VectorData) MeasurePath(segmentName string) (float64, error) {
	item, exists := vd.mapItems[segmentName]
	if !exists {
		return 0, fmt.Errorf("unknown map item '%s'", segmentName)
	}
	var distance float64
	var err error
	switch item := item.(type) {
	case *map_locationType:
		if item.itemType != mitPath {
			return 0, item.Error("%s is not a path", segmentName)
		}
		distance = great.MetersInPath(item.location)
	case *mapSegmentType:
		distance, err = segmentDistance(item)
	case *mapRouteType:
		distance, err = routeDistance(item)
	default:
		return 0, item.Error("%s is not a path, segment, or route", segmentName)
	}
	return distance, err
}

func segmentDistance(seg *mapSegmentType) (float64, error) {
	paths, err := seg.gatherPaths()
	if err != nil {
		return 0, err
	}
	return gatheredPathDistance(paths), nil
}

func routeDistance(route *mapRouteType) (float64, error) {
	paths, err := route.gatherPaths()
	if err != nil {
		return 0, err
	}
	return gatheredPathDistance(paths), nil
}

func gatheredPathDistance(paths []pathLocationInfo) float64 {
	var distance float64
	for _, p := range paths {
		distance += great.MetersInPath(p.path.location)
	}
	return distance
}


func (vd *VectorData) MeasurePathUpTo(segmentName string, upToDistance float64) (foundLat float64,
		foundLong float64, distance float64, pathName string, index int, err error) {
	item, exists := vd.mapItems[segmentName]
	if !exists {
		err = fmt.Errorf("unknown map item '%s'", segmentName)
		return
	}
	var paths []pathLocationInfo
	switch item := item.(type) {
	case *map_locationType:
		if item.itemType != mitPath {
			err = item.Error("%s is not a path", segmentName)
			return
		}
		lat, long := item.location.oppositeEndpoint(pathMatchForward)
		paths = []pathLocationInfo{{item, lat, long, true}}
	case *mapSegmentType:
		paths, err = item.gatherPaths()
	case *mapRouteType:
		paths, err = item.gatherPaths()
	default:
		err = item.Error("%s is not a path, segment, or route", segmentName)
	}
	if err != nil {
		return
	}

	for _, rec := range paths {
		pathDistance := upToDistance - distance
		pathDistance, foundLat, foundLong, index = pathDistanceUpTo(rec, pathDistance)
		distance += pathDistance
		if index > -1 {
			pathName = rec.path.Name()
			return
		}
	}
	err = fmt.Errorf("%s '%s' is only %.1f meters (%.2f miles) long",
		typeMapToName[item.ItemType()], item.Name(), distance,
		distance / great.METERS_PER_MILE)
	return
}


func pathDistanceUpTo(rec pathLocationInfo, upTo float64) (float64, float64, float64, int) {
	var distance, prevDistance, prevLat, prevLong float64
	var pos, posIncrement, index, indexIncrement int
	location := rec.path.location
	count := (len(location) >> 1) - 1
	if rec.forward {
		pos, posIncrement, index, indexIncrement = 2, 2, 1, 1
		prevLat, prevLong = location[0], location[1]
	} else {
		pos, posIncrement, index, indexIncrement = len(location) - 4, -2, count - 1, -1
		prevLat, prevLong = location[pos + 2], location[pos + 3]
	}
	prevLat *= great.DEG_TO_RADIANS
	prevLong *= great.DEG_TO_RADIANS
	for count > 0 {
		count--
		lat := location[pos] * great.DEG_TO_RADIANS
		long := location[pos+1] * great.DEG_TO_RADIANS
		meters := great.MetersBetweenPoints(prevLat, prevLong, lat, long)
		distance += meters
		if distance > upTo {
			diff1 := distance - upTo
			diff2 := upTo - prevDistance
			if diff2 < diff1 {
				pos -= posIncrement
				index -= indexIncrement
				distance = prevDistance
			}
			return distance, location[pos], location[pos+1], index
		}
		pos += posIncrement
		index += indexIncrement
		prevLat, prevLong = lat, long
		prevDistance = distance
	}
	return distance, 0, 0, -1
}






type pathLocationInfo struct {
	path *map_locationType
	farLat, farLong float64
	forward bool
}

func extendContinuousPath(gatheredPaths []pathLocationInfo, path *map_locationType,
		) ([]pathLocationInfo, error) {
	if len(path.location) < 2 {
		return gatheredPaths, nil
	}
	var prevLat, prevLong float64
	var forward bool
	if len(gatheredPaths) > 0 {
		prev := gatheredPaths[len(gatheredPaths) - 1]
		prevLat, prevLong = prev.farLat, prev.farLong
		match := path.location.pathEndpointMatch(prevLat, prevLong)
		if match == noPathMatch {
			gatheredPaths, match = tryFlippingPath0(gatheredPaths, path)
			if match == noPathMatch {
				return nil, path.Error("'%s' does not share an endpoint with '%s'",
					path.Name(), prev.path.Name())
			}
		}
		prevLat, prevLong = path.location.oppositeEndpoint(match)
		forward = match == pathMatchForward
	} else {
		prevLat, prevLong = path.location.oppositeEndpoint(pathMatchForward)
		forward = true
	}
	gatheredPaths = append(gatheredPaths, pathLocationInfo{path, prevLat, prevLong, forward})
	return gatheredPaths, nil
}

func tryFlippingPath0(gatheredPaths []pathLocationInfo, path *map_locationType,
		) ([]pathLocationInfo, int) {
	if len(gatheredPaths) != 1 {
		return nil, noPathMatch
	}
	rec0 := gatheredPaths[0]
	match := pathMatchForward
	if rec0.forward {
		match = pathMatchReverse
	}
	prevLat, prevLong := rec0.path.location.oppositeEndpoint(match)
	match = path.location.pathEndpointMatch(prevLat, prevLong)
	if match == noPathMatch {
		return nil, noPathMatch
	}
	return []pathLocationInfo{{rec0.path, prevLat, prevLong, !rec0.forward}}, match
}


func (ms *mapSegmentType) gatherPaths() ([]pathLocationInfo, error) {
	var gatheredPaths []pathLocationInfo
	var err error
	for _, item := range ms.paths {
		switch item := item.(type) {
		case *map_locationType:
			gatheredPaths, err = extendContinuousPath(gatheredPaths, item)
			if err != nil {
				return nil, err
			}
		case *map_referenceAggregateType:
			for _, mem := range item.targets {
				path, is := mem.(*map_locationType)
				if !is || (path.itemType != mitPath && len(path.location) != 2) {
					return nil, path.Error(
						"aggregate member '%s' is not a path or waypoint",
						path.Name())
				}
				gatheredPaths, err = extendContinuousPath(gatheredPaths, path)
				if err != nil {
					return nil, fmt.Errorf("%s in segment %s", err,
						item.parentName)
				}
			}
		default:
			return nil, item.Error("'%s' is not a path or waypoint", item.Name())
		}
	}
	return gatheredPaths, nil
}


func (route *mapRouteType) gatherPaths() ([]pathLocationInfo, error) {
	var paths [][]pathLocationInfo
	for _, seg := range route.segments {
		switch seg.ItemType() {
		case mitPath, mitPoint, mitMarker, mitCircle:
			item := seg.(*map_locationType)
			lat, long := item.location.oppositeEndpoint(pathMatchForward)
			paths = append(paths, []pathLocationInfo{{item, lat, long, true}})
		case mitSegment:
			item := seg.(*mapSegmentType)
			pts, err := item.gatherPaths()
			if err != nil {
				return nil, err
			}
			paths = append(paths, pts)
		case mitSegments:
			item := seg.(*map_referenceAggregateType)
			for _, mem := range item.targets {
				switch mem.ItemType() {
				case mitSegment:
					sg := mem.(*mapSegmentType)
					pts, err := sg.gatherPaths()
					if err != nil {
						return nil, err
					}
					paths = append(paths, pts)
				case mitPath, mitPoint, mitMarker, mitCircle:
					pth := mem.(*map_locationType)
					lat, long := pth.location.oppositeEndpoint(pathMatchForward)
					paths = append(paths, []pathLocationInfo{{pth, lat, long, true}})
				default:
					return nil, mem.Error("'%s' is not a segment, path, or waypoint",
						mem.Name())
				}
			}
		default:
			return nil, seg.Error("'%s' is not a path, segment, segment list, or waypoint",
				seg.Name())
		}
	}
	var gatheredPaths []pathLocationInfo
	var prevLat, prevLong float64
	for groupX, group := range paths {
		if len(gatheredPaths) > 0 {
			path0 := group[0].path
			pathN := group[len(group)-1].path
			path0Lat, path0Long := path0.location.oppositeEndpoint(pathMatchReverse)
			pathNLat, pathNLong := pathN.location.oppositeEndpoint(pathMatchForward)
			summation := locationPairs{path0Lat, path0Long, pathNLat, pathNLong}
			match := summation.pathEndpointMatch(prevLat, prevLong)
			if match == noPathMatch {
				// Try flipping the segment
				path0Lat, path0Long := path0.location.oppositeEndpoint(
					pathMatchForward)
				pathNLat, pathNLong := pathN.location.oppositeEndpoint(
					pathMatchReverse)
				summation := locationPairs{path0Lat, path0Long, pathNLat, pathNLong}
				match = summation.pathEndpointMatch(prevLat, prevLong)
			}
			if match == noPathMatch {
				gatheredPaths, match = tryFlippingGroup0(groupX, gatheredPaths,
					group)
				if match == noPathMatch {
					return nil, path0.Error(
						"path '%s' does not join with route '%s'",
						path0.Name(), route.Name())
				}
			}
			if match == pathMatchReverse {
				reversePathGroup(group)
			}
		}
		gatheredPaths = append(gatheredPaths, group...)
		pathN := group[len(group)-1]
		prevLat, prevLong = pathN.farLat, pathN.farLong
	}
	return gatheredPaths, nil
}

func tryFlippingGroup0(groupX int, gatheredPaths []pathLocationInfo, group []pathLocationInfo,
		) ([]pathLocationInfo, int) {
	if groupX > 1 {
		return nil, noPathMatch
	}
	path0 := group[0].path
	pathN := group[len(group)-1].path
	gatheredLocation := gatheredPaths[0].path.location
	var match int
	for _, pointX := range []int{0, len(gatheredLocation) - 2} {
		prevLat, prevLong := gatheredLocation[pointX], gatheredLocation[pointX + 1]
		path0Lat, path0Long := path0.location.oppositeEndpoint(pathMatchReverse)
		pathNLat, pathNLong := pathN.location.oppositeEndpoint(pathMatchForward)
		summation := locationPairs{path0Lat, path0Long, pathNLat, pathNLong}
		match = summation.pathEndpointMatch(prevLat, prevLong)
		if match != noPathMatch {
			break
		}
		path0Lat, path0Long = path0.location.oppositeEndpoint(pathMatchForward)
		pathNLat, pathNLong = pathN.location.oppositeEndpoint(pathMatchReverse)
		summation = locationPairs{path0Lat, path0Long, pathNLat, pathNLong}
		match = summation.pathEndpointMatch(prevLat, prevLong)
		if match != noPathMatch {
			break
		}
	}
	if match == noPathMatch {
		return nil, noPathMatch
	}
	reversePathGroup(gatheredPaths)
	return gatheredPaths, match
}


func reversePathGroup(group []pathLocationInfo) {
	for i, j := 0, len(group) - 1; i < j; i, j = i+1, j-1 {
		group[i], group[j] = group[j], group[i]
	}
	for i, pth := range group {
		match := pathMatchForward
		if pth.forward {
			match = pathMatchReverse
		}
		pth.farLat, pth.farLong = pth.path.location.oppositeEndpoint(match)
		pth.forward = !pth.forward
		group[i] = pth
	}
}

