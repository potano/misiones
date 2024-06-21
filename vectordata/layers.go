// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import ( 
        "potano.misiones/sexp"
)


type mapLayersType struct {
	mapItemCore
	layers []*mapLayerType
}

func newMapLayers(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	obj := &mapLayersType{}
	obj.source = source
	doc.layersRoot = obj
	name, err := doc.registerMapItem(obj, listName)
	obj.name = name
	obj.itemType = mitLayers
	return obj, err
}

func (ml *mapLayersType) addFeature(layer mapItemType) {
	ml.layers = append(ml.layers, layer.(*mapLayerType))
}



type mapLayerType struct {
	mapItemCore
	menuitem string
	features []mapItemType
}

func newMapLayer(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ml := &mapLayerType{}
	ml.source = source
	name, err := doc.registerMapItem(ml, listName)
	ml.name = name
	ml.itemType = mitLayer
	return ml, err
}

func (ml *mapLayerType) setMenuitem(menuitem *map_textType) {
	ml.menuitem = menuitem.text
}

func (ml *mapLayerType) addFeature(feature mapItemType) {
	ml.features = append(ml.features, feature)
}

