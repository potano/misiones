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
	err = fmt.Errorf("%s %s is only %.1f meters (%.2f miles) long",
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
		pos, posIncrement, index, indexIncrement = len(location) - 4, -2, count - 2, -1
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
			return nil, path.Error("'%s' does not share and endpoint with '%s'",
				path.Name(), prev.path.Name())
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
				if !is || path.itemType != mitPath {
					return nil, path.Error(
						"aggregate member '%s' is not a path", path.Name())
				}
				gatheredPaths, err = extendContinuousPath(gatheredPaths, path)
				if err != nil {
					return nil, path.Error("in segment %s, %s",
						item.parentName, err)
				}
			}
		default:
			return nil, item.Error("'%s' is not a path", item.Name())
		}
	}
	return gatheredPaths, nil
}


func (route *mapRouteType) gatherPaths() ([]pathLocationInfo, error) {
	var paths [][]pathLocationInfo
	for _, seg := range route.segments {
		switch item := seg.(type) {
		case *map_locationType:
			if item.itemType != mitPath {
				return nil, item.Error("'%s' is not a path", item.Name())
			}
			lat, long := item.location.oppositeEndpoint(pathMatchForward)
			paths = append(paths, []pathLocationInfo{{item, lat, long, true}})
		case *mapSegmentType:
			pts, err := item.gatherPaths()
			if err != nil {
				return nil, err
			}
			paths = append(paths, pts)
		case *map_referenceAggregateType:
			for _, mem := range item.targets {
				sg, is := mem.(*mapSegmentType)
				if !is {
					return nil, mem.Error("'%s' is not a segment", mem.Name())
				}
				pts, err := sg.gatherPaths()
				if err != nil {
					return nil, err
				}
				paths = append(paths, pts)
			}
		default:
			return nil, seg.Error("'%s' is not a path, segment, or segment list",
				seg.Name())
		}
	}
	var gatheredPaths []pathLocationInfo
	var prevLat, prevLong float64
	for _, group := range paths {
		if len(gatheredPaths) > 0 {
			path0 := group[0].path
			summation := locationPairs{path0.location[0], path0.location[1],
				group[0].farLat, group[0].farLong}
			match := summation.pathEndpointMatch(prevLat, prevLong)
			if match == noPathMatch {
				return nil, path0.Error(
					"path '%s' does not join with route '%s'",
					path0.Name(), route.Name())
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

