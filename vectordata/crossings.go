// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"sort"
)

type locationIndexType int32

// Records passed as input to the crossing finder
type locationPathsRecord struct {
	locationIndex, startOffset locationIndexType
	pairs locationPairs
}

// Records for crossing points along paths and segments
// The crossing finder generates the initial sets of these for each path
type latlongRef struct {
	point latlongType
	indices[3] locationIndexType
}

type latlongRefs []latlongRef



type crossingFinderType struct {
	crossingFinderChannel chan locationPathsRecord
	crossingsDoneChannel chan map[locationIndexType]latlongRefs
}

func newCrossingFinder() *crossingFinderType {
	cf := &crossingFinderType{
		crossingFinderChannel: make(chan locationPathsRecord, 100),
		crossingsDoneChannel: make(chan map[locationIndexType]latlongRefs, 1),
	}

	go findCrosspoints(cf)
	return cf
}

func (cf *crossingFinderType) addLocation(locationIndex, startIndex locationIndexType,
		pairs locationPairs) {
	cf.crossingFinderChannel <- locationPathsRecord{locationIndex, startIndex, pairs}
}

func (cf *crossingFinderType) signalNoMoreInput() {
	close(cf.crossingFinderChannel)
}

func (cf *crossingFinderType) getAllCrosspoints() map[locationIndexType]latlongRefs {
	return <- cf.crossingsDoneChannel
}



type cpathInfo struct {
	locationIndex, offset locationIndexType
}

func findCrosspoints(cf *crossingFinderType) {
	allPoints := map[latlongType][]cpathInfo{}
	for {
		locpath, ok := <- cf.crossingFinderChannel
		if !ok {
			break
		}
		lastIndex := len(locpath.pairs) - 2
		for i := 0; i <= lastIndex; i += 2 {
			pair := locpath.pairs.latlongPair(i)
			newCpi := cpathInfo{locpath.locationIndex,
				locationIndexType(i) + locpath.startOffset}
			cpiList := allPoints[pair]
			if cpiList == nil {
				cpiList = []cpathInfo{newCpi}
			} else {
				cpiList = append(cpiList, newCpi)
			}
			if i == 0 || i == lastIndex {
				// Ensure that path endpoints are always treated as crosspoints
				cpiList = append(cpiList, cpathInfo{-1, 0})
			}
			allPoints[pair] = cpiList
		}
	}

	locationCrosspoints := map[locationIndexType]latlongRefs{}

	for point, cpiList := range allPoints {
		if len(cpiList) > 1 {
			for _, cp := range cpiList {
				locationCrosspoints[cp.locationIndex] =
					append(locationCrosspoints[cp.locationIndex],
					latlongRef{point, [3]locationIndexType{cp.offset, 0, 0}})
			}
			for _, refs := range locationCrosspoints {
				refs.sort()
			}
		}
	}

	cf.crossingsDoneChannel <- locationCrosspoints
}







func (ll1 latlongRef) comesBefore(other latlongRef) bool {
	return ll1.indices[0] < other.indices[0] ||
		(ll1.indices[0] == other.indices[0] &&
			(ll1.indices[1] < other.indices[1] ||
				(ll1.indices[1] == other.indices[1] &&
					ll1.indices[2] < other.indices[2])))
}

func (llr latlongRef) clone() (newRef latlongRef) {
	newRef.point = llr.point
	newRef.indices[0] = llr.indices[0]
	newRef.indices[1] = llr.indices[1]
	newRef.indices[2] = llr.indices[2]
	return newRef
}

func (llr latlongRef) cloneAndPushLevel(ind locationIndexType) (newRef latlongRef) {
	newRef.point = llr.point
	newRef.indices[2] = llr.indices[1]
	newRef.indices[1] = llr.indices[0]
	newRef.indices[0] = ind
	return newRef
}


func (llr latlongRefs) findMatchingCrosspoints(others latlongRefs,
		) (match, nomatch latlongRefs, ambiguous bool) {
	var found, isSamePoint bool
	for _, info := range llr {
		found = false
		for _, otherInfo := range others {
			if info.point.samePoint(otherInfo.point) {
				found, isSamePoint = true, true
				for _, pt := range match {
					if !info.point.samePoint(pt.point) {
						ambiguous = true
						isSamePoint = false
						break
					}
				}
				break
			}
		}
		if !found || !isSamePoint {
			nomatch = append(nomatch, info)
		} else if isSamePoint {
			match = append(match, info)
			// we deal later with the existence of multiple matching paths
		}
	}
	return
}

func (llr latlongRef) String() string {
	return fmt.Sprintf("%d [%d,%d,%d]", llr.point, llr.indices[0], llr.indices[1],
		llr.indices[2])
}

func (llrs latlongRefs) strings() []string {
	out := make([]string, len(llrs))
	for i, r := range llrs {
		out[i] = r.String()
	}
	return out
}

func (llrs latlongRefs) sort() {
	if len(llrs) > 1 {
		sort.Slice(llrs, func (i, j int) bool {
			a, b := llrs[i].indices[0], llrs[j].indices[0]
			if a != b {
				return a < b
			}
			a, b = llrs[i].indices[1], llrs[j].indices[1]
			if a != b {
				return a < b
			}
			return llrs[i].indices[2] < llrs[j].indices[2]
		})
	}
}

// For cases when the crossing finder cannot find a match
func synthesizeCrosspointReferences(item threadableMapItemType) (startRefs, endRefs latlongRefs) {
	startRefs = synthesizeCrosspointEndReference(item, true)
	endRefs = synthesizeCrosspointEndReference(item, false)
	return
}

func synthesizeCrosspointEndReference(item threadableMapItemType, makeStartpoint bool) latlongRefs {
	altPoint, usePoint, altOffset, useOffset := item.endpointsAndOffsets()
	if makeStartpoint {
		usePoint, useOffset = altPoint, altOffset
	}
	ref := latlongRef{usePoint, [3]locationIndexType{useOffset, 0, 0}}
	depth := locationIndexType(1)
	for {
		if rc, is := item.(*threadableMapItemReference); is {
			item = rc.item
		}
		if rs, is := item.(*mapRouteOrSegmentType); !is {
			break
		} else {
			item = rs.children[useOffset].(threadableMapItemType)
		}
		itemStartPoint, _, altOffset, useOffset := item.endpointsAndOffsets()
		if itemStartPoint == usePoint {
			useOffset = altOffset
		}
		ref.indices[depth] = useOffset
		depth++
	}
	return latlongRefs{ref}
}

