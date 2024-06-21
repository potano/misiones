// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

// Establishes the exact traversal of a threadable map item (a route or a segment) via its
// constituent parts.  In the normal case, each constituent shares a point of intersection (i.e.
// a latitude/longitude pair) with its neighbor to either side.  The resulting threaded route or
// segment includes only those points which fall between the intersection points of the respective
// constituent.
func (vd *VectorData) threadRouteOrSegment(item *mapRouteOrSegmentType) error {
	children, err := gatherThreadedItemList(item)
	if nil != err {
		return err
	}
	markedChildren, err := vd.markComponentIntersections(item, children)
	if err != nil {
		return err
	}
	pickedChildren := pickThreadedItems(item, markedChildren)
	return vd.finishThreading(pickedChildren)
}


// Special methods for path, point, point-like, route, and segment items
type threadableMapItemType interface {
	mapItemType
	getCrosspoints() latlongRefs
	endpointsAndOffsets() (pt1, pt2 latlongType, off1, off2 locationIndexType)
	oppositeEndpoint(startPoint latlongType) (endPoint latlongType, endOffset,
		startOffset locationIndexType)
	isPoint() bool
	resolveReferenceToLocation(ref latlongRef) *map_locationType
}


// Markup for each child of a item being threaded
type pendingChildInfo struct {
	child threadableMapItemType
	startRefs, endRefs latlongRefs
}



func gatherThreadedItemList(item *mapRouteOrSegmentType) ([]threadableMapItemType, error) {
	var children []threadableMapItemType
	for _, listItem := range item.routeComponents() {
		switch listItem := listItem.(type) {
		case *map_locationType, *mapRouteOrSegmentType:
			children = append(children, listItem.(threadableMapItemType))
		case *map_referenceAggregateType:
			for targX, refTarget := range listItem.targets {
				switch refTarget := refTarget.(type) {
				case *map_locationType, *mapRouteOrSegmentType:
					children = append(children,
						&threadableMapItemReference{
							refTarget.(threadableMapItemType),
							listItem.names[targX].Source()})
				default:
					return nil, refTarget.Error("'%s' not allowed in %s",
						refTarget.Name(), listItem.Name())
				}
			}
		default:
			return nil, item.Error("'%s' not allowed in %s", listItem.Name(),
				item.Name())
		}
	}
	return children, nil
}


func (vd *VectorData) markComponentIntersections(item threadableMapItemType,
		children []threadableMapItemType) ([]pendingChildInfo, error) {
	// Augment each component item with a record that indicates how it relates to its neighbor
	boundedChildren := make([]pendingChildInfo, len(children))
	previousChildInfo := &boundedChildren[0]
	var previousCrosspoints latlongRefs
	endingCrosspoints := make([]latlongRefs, len(children))
	for childX, child := range children {
		childInfo := &boundedChildren[childX]
		childInfo.child = child
		childCrosspoints := child.getCrosspoints()
		if childX > 0 {
			matching, notMatching, ambiguous :=
				childCrosspoints.findMatchingCrosspoints(previousCrosspoints)
			if len(matching) > 0 {
				childInfo.startRefs = matching
				if ambiguous {
					noteAmbiguousCrosspointMatch(vd, child,
					children[childX - 1])
				}
				matching, _, _ = previousChildInfo.child.getCrosspoints().
					findMatchingCrosspoints(matching)
				previousChildInfo.endRefs = matching
			}
			if child.isPoint() {
				childInfo.endRefs = childInfo.startRefs
				previousCrosspoints = matching
			} else {
				previousCrosspoints = notMatching
			}
		} else {
			previousCrosspoints = childCrosspoints
		}
		endingCrosspoints[childX] = previousCrosspoints
		previousChildInfo = childInfo
	}

	// Flag points of discontinuity and fill in missing endpoints
	for childX := range boundedChildren {
		childInfo := &boundedChildren[childX]
		startRefs, endRefs := childInfo.startRefs, childInfo.endRefs
		child := childInfo.child
		if len(startRefs) == 0 || len(endRefs) == 0 {
			childEnds1, childEnds2 := synthesizeCrosspointReferences(child)
			crosspointsAtEnd := endingCrosspoints[childX]
			availableForStartpoint := crosspointsAtEnd
			if len(startRefs) == 0 && len(endRefs) == 0 {
				if len(boundedChildren) > 1 {
					noteFailedCrosspointMatch(vd, child, item)
				}
				childInfo.startRefs = childEnds1
				childInfo.endRefs = childEnds2
				continue
			} else if len(endRefs) == 0 {
				matching1, notMatching1, _ := childEnds1.
					findMatchingCrosspoints(crosspointsAtEnd)
				matching2, notMatching2, _ := childEnds2.
					findMatchingCrosspoints(crosspointsAtEnd)
				if len(matching2) > 0 {
					if len(matching1) > 0 {
						noteAmbiguousEndpoint(vd, child, item)
					}
					endRefs = matching2
					availableForStartpoint = notMatching2
				} else if len(matching1) > 0 {
					endRefs = matching1
					availableForStartpoint = notMatching1
				} else {
					endRefs = childEnds2
				}
			} else {
				_, availableForStartpoint, _ = crosspointsAtEnd.
					findMatchingCrosspoints(endRefs)
			}
			if len(startRefs) == 0 {
				if childX > 0 {
					noteFailedCrosspointMatch(vd, child, item)
				}
				matching1, _, _ := childEnds1.
					findMatchingCrosspoints(availableForStartpoint)
				matching2, _, _ := childEnds2.
					findMatchingCrosspoints(availableForStartpoint)
				if len(matching1) > 0 {
					if len(matching2) > 0 {
						noteAmbiguousEndpoint(vd, child, item)
					}
					startRefs = matching1
				} else if len(matching2) > 0 {
					startRefs = matching2
				} else {
					startRefs = childEnds1
				}
			}
		}
		if len(startRefs) > 1 || len(endRefs) > 1 {
			ascending := startRefs[0].comesBefore(endRefs[0])
			if len(startRefs) > 1 {
				startRefs = resolveEndRefs(child, startRefs, ascending)
			}
			if len(endRefs) > 1 {
				endRefs = resolveEndRefs(child, endRefs, !ascending)
			}
		}
		childInfo.startRefs = startRefs
		childInfo.endRefs = endRefs
	}
	for childX, child := range boundedChildren {
		if item, is := child.child.(*threadableMapItemReference); is {
			boundedChildren[childX].child = item.item
		}
	}
	return boundedChildren, nil
}


func resolveEndRefs(child threadableMapItemType, refs latlongRefs, reverse bool) latlongRefs {
	pos, inc, stop := 0, 1, len(refs)
	if reverse {
		pos, inc, stop = len(refs) - 1, -1, -1
	}
	havePath := false
	var takePos int
	for pos != stop {
		ref := refs[pos]
		loc := child.resolveReferenceToLocation(ref)
		if loc.isPoint() {
			takePos = pos
		} else if havePath {
			break
		} else {
			havePath = true
			takePos = pos
		}
		pos += inc
	}
	return latlongRefs{refs[takePos]}
}



func noteAmbiguousCrosspointMatch(vd *VectorData, child, parent mapItemType) {
	logThreadingError(vd, child, parent, "%s %s connects with %s %s at multiple points")
}

func noteAmbiguousEndpoint(vd *VectorData, child, parent mapItemType) {
	logThreadingError(vd, child, parent, "cannot determine free endpoint of %s %s under %s %s")
}

func noteFailedCrosspointMatch(vd *VectorData, child, parent mapItemType) {
	logThreadingError(vd, child, parent, "%s %s does not connect with %s %s")
}



type pickedItem struct {
	item threadableMapItemType			// route, segment, path, marker, or point
	children []pickedItem				// selected children of segment
	startPoint, endPoint latlongType		// latitude and longitude of endpoints
	startOffset, endOffset locationIndexType	// start and end offsets within path
	shortened bool					// TRUE if we need to generate new item
}

// Returns a picked-item structure.  The route or segment at the root of this structure is the
// route or segment being threaded along with information about its endpoints.  The children of
// this root element are trees referenced in the definition of the root element such that each
// tree contains only those paths and segments which lie between the indicated intersection
// points.  Later parts of the processing reduces these reference subtrees to a form suitable
// for the root element.
func pickThreadedItems(root mapItemType, marked []pendingChildInfo) pickedItem {
	if len(marked) == 0 {
		return pickedItem{item: root.(threadableMapItemType)}
	}
	pickedItems := make([]pickedItem, 0, len(marked))
	for _, pending := range marked {
		child := pending.child
		alignPoint, endPoint, startOffset, endOffset := child.endpointsAndOffsets()
		crossStartOffset := pending.startRefs[0].indices[0]
		crossEndOffset := pending.endRefs[0].indices[0]
		if pending.endRefs[0].comesBefore(pending.startRefs[0]) {
			alignPoint, endPoint = endPoint, alignPoint
			startOffset, endOffset = endOffset, startOffset
		}
		picked := pending.pickItem(child, 1, alignPoint, endPoint,
			startOffset, endOffset, crossStartOffset, crossEndOffset)
		pickedItems = append(pickedItems, picked)
	}
	rootStartOffset := locationIndexType(0)
	rootEndOffset := locationIndexType(len(pickedItems) - 1)
	rootStartPoint := pickedItems[0].startPoint
	rootEndPoint := pickedItems[rootEndOffset].endPoint
	return pickedItem{
		item: root.(threadableMapItemType),
		children: pickedItems,
		startPoint: rootStartPoint,
		endPoint: rootEndPoint,
		startOffset: rootStartOffset,
		endOffset: rootEndOffset,
	}
}



func (pc *pendingChildInfo) pickItem(item mapItemType, depth int, alignPoint, endPoint latlongType,
		alignPointOffset, endPointOffset locationIndexType,
		crossStartOffset, crossEndOffset locationIndexType) pickedItem {
	selectionStart, selectionEnd := alignPointOffset, endPointOffset
	if crossStartOffset >= 0 {
		selectionStart = crossStartOffset
	}
	if crossEndOffset >= 0 {
		selectionEnd = crossEndOffset
	}
	partial := (selectionStart != alignPointOffset || selectionEnd != endPointOffset) &&
		(selectionStart != endPointOffset || selectionEnd != alignPointOffset)
	if path, is := item.(*map_locationType); is {
		return pickedItem{
			item: item.(threadableMapItemType),
			startPoint: path.pointAtOffset(selectionStart),
			endPoint: path.pointAtOffset(selectionEnd),
			startOffset: selectionStart,
			endOffset: selectionEnd,
			shortened: partial,
		}
	}

	readPos := alignPointOffset
	var readIncrement locationIndexType
	descendingAlignment := endPointOffset < alignPointOffset
	if descendingAlignment {
		readIncrement = -1
	} else {
		readIncrement = 1
	}
	descendingSelection := (selectionEnd < selectionStart) != descendingAlignment
	minSelection, maxSelection := selectionStart, selectionEnd
	if maxSelection < minSelection {
		minSelection, maxSelection = maxSelection, minSelection
	}
	var skipCount int
	if descendingAlignment {
		skipCount = int(maxSelection - alignPointOffset)
	} else {
		skipCount = int(alignPointOffset - minSelection)
	}
	if skipCount < 0 {
		skipCount = - skipCount
	}
	pickCount := int(1 + maxSelection - minSelection)
	var writePos, writeIncrement int
	if descendingSelection {
		writePos = pickCount - 1
		writeIncrement = -1
	} else {
		writePos = 0
		writeIncrement = 1
	}
	children := item.(*mapRouteOrSegmentType).children
	pickedChildren := make([]pickedItem, pickCount)
	var parentStartPoint, parentEndPoint latlongType
	for {
		child := children[readPos]
		farEndpoint, farEndpointOffset, nearEndpointOffset :=
			child.(threadableMapItemType).oppositeEndpoint(alignPoint)
		if skipCount > 0 {
			skipCount--
			readPos += readIncrement
			alignPoint = farEndpoint
			continue
		}
		var childCrossingStart, childCrossingEnd locationIndexType
		if readPos == crossStartOffset {
			childCrossingStart = pc.startRefs[0].indices[depth]
		} else {
			childCrossingStart = -1
		}
		if readPos == crossEndOffset {
			childCrossingEnd = pc.endRefs[0].indices[depth]
		} else {
			childCrossingEnd = -1
		}
		picked := pc.pickItem(child, depth + 1, alignPoint, farEndpoint,
			nearEndpointOffset, farEndpointOffset, childCrossingStart, childCrossingEnd)
		pickedChildren[writePos] = picked
		writePos += writeIncrement
		partial = partial || picked.shortened
		if readPos == selectionStart {
			parentStartPoint = picked.startPoint
		}
		if readPos == selectionEnd {
			parentEndPoint = picked.endPoint
		}
		if pickCount < 2 {
			break
		}
		alignPoint = farEndpoint
		pickCount--
		readPos += readIncrement
	}
	return pickedItem{
		item: item.(threadableMapItemType),
		children: pickedChildren,
		startPoint: parentStartPoint,
		endPoint: parentEndPoint,
		startOffset: 0,
		endOffset: locationIndexType(len(pickedChildren) - 1),
		shortened: partial,
	}
}



func (vd *VectorData) finishThreading(rootPickedItem pickedItem) error {
	item := rootPickedItem.item.(*mapRouteOrSegmentType)
	item.setEndpointsAndChildren(rootPickedItem.startPoint, rootPickedItem.endPoint,
		rootPickedItem.fixChildren(vd, item, item.ItemType() == mitRoute))
	return nil
}


func (pi pickedItem) fixChildren(vd *VectorData, parent mapItemType,
		returnSegments bool) []mapItemType {
	var readyChildren []mapItemType
	for _, picked := range pi.children {
		switch child := picked.item.(type) {
		case *mapRouteOrSegmentType:
			if child.ItemType() == mitRoute || !returnSegments {
				grandchildren := picked.fixChildren(vd, parent, returnSegments)
				readyChildren = append(readyChildren, grandchildren...)
				continue
			}
			if picked.shortened {
				child = child.clone(vd, parent)
				child.setEndpointsAndChildren(picked.startPoint, picked.endPoint,
					picked.fixChildren(vd, child, returnSegments))
			}
			readyChildren = append(readyChildren, child)
		case *map_locationType:
			if picked.shortened {
				child = child.makeSubpath(vd, parent,
					picked.startOffset, picked.endOffset)
			}
			readyChildren = append(readyChildren, child)
		}
	}
	return readyChildren
}

