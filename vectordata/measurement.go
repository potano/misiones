// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"

	"potano.misiones/great"
)

type measurementWalker interface {
	measurePath(path *map_locationType, startOffset, endOffset locationIndexType) bool
}



func (vd *VectorData) MeasurePath(itemName string) (float64, error) {
	measurer := &simplePathMeasurer{}
	err := vd.walkPathsForNamedItem(measurer, itemName, false)
	if err != nil {
		return 0, err
	}
	return measurer.meters, nil
}

type simplePathMeasurer struct {
	meters float64
}

func (spm *simplePathMeasurer) measurePath(path *map_locationType,
		startOffset, endOffset locationIndexType) bool {
	spm.meters += great.MetersInPath(path.location.asFloatSlice())
	return true
}


func (vd *VectorData) MeasurePathUpTo(itemName string, upToDistance float64) (foundLat float64,
		foundLong float64, distance float64, pathName string, index int, err error) {
	measurer := &upToDistanceMeasurer{upToDistance: upToDistance, index: -1}
	err = vd.walkPathsForNamedItem(measurer, itemName, false)
	if err != nil {
		index = -1
		return
	}
	foundLat, foundLong = measurer.lat, measurer.long
	distance, pathName, index = measurer.distance, measurer.pathName, measurer.index
	if index < 0 {
		item := vd.mapItems[itemName]
		err = fmt.Errorf("%s '%s' is only %.1f meters (%.2f miles) long",
			item.ItemTypeString(), item.Name(), distance,
			distance / great.METERS_PER_MILE)
	}
	return
}

type upToDistanceMeasurer struct {
	lat, long, upToDistance, distance float64
	pathName string
	index int
}

func (updm *upToDistanceMeasurer) measurePath(path *map_locationType,
		startOffset, endOffset locationIndexType) bool {
	reverse := endOffset < startOffset
	var pairs []float64
	if reverse {
		pairs = path.location.asReverseFloatSlice()
	} else {
		pairs = path.location.asFloatSlice()
	}
	steps := great.MetersBetweenPointPairs(pairs)
	distance := updm.distance
	for index, step := range steps {
		if distance + step >= updm.upToDistance {
			diff1 := distance + step - updm.upToDistance
			diff2 := updm.upToDistance - distance
			if diff1 < diff2 {
				index++
				distance += step
			}
			updm.lat = pairs[index*2]
			updm.long = pairs[index*2 + 1]
			updm.distance = distance
			updm.pathName = path.Name()
			if reverse {
				updm.index = len(steps) - index
			} else {
				updm.index = index
			}
			return false
		}
		distance += step
	}
	updm.distance = distance
	return true
}





func (vd *VectorData) walkPathsForNamedItem(walker measurementWalker, name string,
		reverse bool) error {
	item, exists := vd.mapItems[name]
	if !exists {
		return fmt.Errorf("unknown map item '%s'", name)
	}
	var align, endpoint latlongType
	var startOffset, endOffset locationIndexType
	if tItem, is := item.(threadableMapItemType); !is {
		return fmt.Errorf("%s %s is not a threadable type", item.ItemTypeString(),
			item.Name())
	} else {
		align, endpoint, startOffset, endOffset = tItem.endpointsAndOffsets()
	}
	if reverse {
		align = endpoint
		startOffset, endOffset = endOffset, startOffset
	}
	switch tItem := item.(type) {
	case *map_locationType:
		walker.measurePath(tItem, startOffset, endOffset)
	case *mapRouteOrSegmentType:
		walkRouteOrSegment(walker, tItem, align, startOffset, endOffset)
	}
	return nil
}

func walkRouteOrSegment(walker measurementWalker, item *mapRouteOrSegmentType, align latlongType,
		pos, endOffset locationIndexType) bool {
	children := item.routeComponents()
	increment := locationIndexType(1)
	if endOffset < pos {
		increment = -1
	}
	for {
		child := children[pos].(threadableMapItemType)
		oppositeEndpoint, oppositeOffset, nearOffset := child.oppositeEndpoint(align)
		switch tChild := child.(type) {
		case *map_locationType:
			if !walker.measurePath(tChild, nearOffset, oppositeOffset) {
				return false
			}
		case *mapRouteOrSegmentType:
			if !walkRouteOrSegment(walker, tChild, align, nearOffset, oppositeOffset) {
				return false
			}
		}
		if pos == endOffset {
			break
		}
		align = oppositeEndpoint
		pos += increment
	}
	return true
}

