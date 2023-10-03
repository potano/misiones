// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"

	"potano.misiones/great"
)


func (vd *VectorData) MeasurePath(segmentName string) (float64, error) {
	segs, err := vd.gatherSegmentsForNamedItem(segmentName)
	if err != nil {
		return 0, err
	}
	var distance float64
	for _, segment := range segs {
		for _, path := range segment.paths {
			points, _, _ := path.points()
			distance += great.MetersInPath(points)
		}
	}
	return distance, nil
}


func (vd *VectorData) MeasurePathUpTo(segmentName string, upToDistance float64) (foundLat float64,
		foundLong float64, distance float64, pathName string, index int, err error) {
	segs, err := vd.gatherSegmentsForNamedItem(segmentName)
	if err != nil {
		return
	}
	for _, segment := range segs {
		for _, rec := range segment.paths {
			pathDistance := upToDistance - distance
			pathDistance, foundLat, foundLong, index =
				pathDistanceUpTo(rec, pathDistance)
			distance += pathDistance
			if index > -1 {
				pathName = rec.path.Name()
				return
			}
		}
	}
	item := vd.mapItems[segmentName]
	err = fmt.Errorf("%s '%s' is only %.1f meters (%.2f miles) long",
		typeMapToName[item.ItemType()], item.Name(), distance,
		distance / great.METERS_PER_MILE)
	return
}



func (vd *VectorData) gatherSegmentsForNamedItem(name string) ([]*gatheredSegment, error) {
	item, exists := vd.mapItems[name]
	if !exists {
		return nil, fmt.Errorf("unknown map item '%s'", name)
	}
	var segments []*gatheredSegment
	var segment *gatheredSegment
	var err error
	switch item := item.(type) {
	case *map_locationType:
		if item.itemType != mitPath {
			return nil, item.Error("%s is not a path", name)
		}
		segment = item.pathAsGatheredSegment()
	case *mapSegmentType:
		segment, err = item.threadPaths()
	case *mapRouteType:
		segments, err = item.threadSegments()
	default:
		return nil, item.Error("%s is not a path, segment, or route", name)
	}
	if segment != nil {
		segments = []*gatheredSegment{segment}
	}
	return segments, err
}


func pathDistanceUpTo(rec gatheredPath, upTo float64) (float64, float64, float64, int) {
	var distance, prevDistance, prevLat, prevLong float64
	var pos, posIncrement, index, indexIncrement int
	location, baseIndex, forward := rec.points()
	count := (len(location) >> 1) - 1
	if forward {
		pos, posIncrement, index, indexIncrement = 0, 2, 1, 1
	} else {
		pos, posIncrement, index, indexIncrement = len(location) - 2, -2, count - 1, -1
	}
	prevLat, prevLong = location[pos], location[pos+1]
	prevLat *= great.DEG_TO_RADIANS
	prevLong *= great.DEG_TO_RADIANS
	for count > 0 {
		count--
		pos += posIncrement
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
			return distance, location[pos], location[pos+1], index + baseIndex
		}
		index += indexIncrement
		prevLat, prevLong = lat, long
		prevDistance = distance
	}
	return distance, 0, 0, -1
}

