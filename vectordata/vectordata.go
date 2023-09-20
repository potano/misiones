// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"sort"
	"potano.misiones/sexp"
)


type VectorData struct {
	mapItems map[string]mapItemType
	inDependencyOrder []string
	referenceItems []*map_referenceAggregateType
	layersRoot *mapLayersType
	styler *styler
	attester *attester
}

type mapItemType interface {
	Name() string
	Source() sexp.ValueSource
	ItemType() int
	noteReferrer(string, mapItemType) error
	Referrers() []string
	addScalars(targetName string, scalars []sexp.LispScalar) error
	addLayer(layer *mapLayerType)
	setMenuitem(layer *map_textType)
	setPopup(popup *mapPopupType)
	setStyle(style *mapStyleType)
	setAttestation(attestation *mapAttestationType)
	setHtml(html *map_textType)
	setRadius(radius *mapRadiusType)
	addFeature(feature mapItemType)
	setConfigurationItem(item mapItemType) error
	styleAndAttestation() (*mapStyleType, *mapAttestationType)
	Error(msg string, args ...any) error
	generateJs() string
}


func NewVectorData() *VectorData {
	return &VectorData{
		mapItems: map[string]mapItemType{},
	}
}

func (vd *VectorData) registerMapItem(item mapItemType, name string) (string, error) {
	if len(name) == 0 {
		name = fmt.Sprintf("$%d", len(vd.mapItems))
	}
	if _, exists := vd.mapItems[name]; exists {
		return name, item.Error("duplicate use of name '%s'", name)
	}
	vd.mapItems[name] = item
	return name, nil
}

func (vd *VectorData) registerReferenceItem(item *map_referenceAggregateType) {
	vd.referenceItems = append(vd.referenceItems, item)
}

func (vd *VectorData) ResolveReferences() error {
	for _, item := range vd.referenceItems {
		err := item.resolveTargets(vd)
		if err != nil {
			return err
		}
	}
	// The parser constructs a DAG from the root element plus zero or more disconnected
	// segments that are also acyclic.  The target-resolution step above aims to join the
	// disconnected segments into the main graph but cannot account for two pathologies:
	// cycles and orphan segments.  We resolve these by examining the inbound edges of
	// the nodes.  Orphan nodes have no inbound edges; the set of inbound edges from each
	// remaining node is used to compute the set of outbound nodes for each node.
	// The cycle-detection analysis starts by identifying the leaf nodes and through
	// successive iterations finds all the nodes which are either leaf nodes or point to
	// nodes which end only at leaf nodes.  Iterations continue until all the nodes in
	// the graph are considered safe or until an interation fails to include any new
	// nodes.  The function reports the graph as containing a cycle if an iteration fails
	// to add any nodes to the safe set.
	childNodesForNode := map[string][]string{}
	for name, node := range vd.mapItems {
		if len(name) > 0 && name[0] != '$' && len(node.Referrers()) == 0 {
			return node.Error("%s '%s' is an orphan",
				typeMapToName[node.ItemType()], node.Name())
		}
		for _, referrer := range node.Referrers() {
			if list, exists := childNodesForNode[referrer]; exists {
				list = append(list, name)
				childNodesForNode[referrer] = list
			} else {
				childNodesForNode[referrer] = []string{name}
			}
		}
	}
	okNodes := map[string]bool{}
	vd.inDependencyOrder = make([]string, len(vd.mapItems))
	for name := range vd.mapItems {
		if _, exists := childNodesForNode[name]; !exists {
			vd.inDependencyOrder[len(okNodes)] = name
			okNodes[name] = true
		}
	}
	for len(okNodes) < len(vd.mapItems) {
		numAdded := 0
		for name, list := range childNodesForNode {
			if !okNodes[name] {
				ok := true
				for _, child := range list {
					if !okNodes[child] {
						ok = false
						break
					}
				}
				if ok {
					vd.inDependencyOrder[len(okNodes)] = name
					okNodes[name] = true
					numAdded++
				}
			}
		}
		if numAdded == 0 {
			return vd.implicateCycleHead(childNodesForNode, okNodes)
		}
	}
	return nil
}


func (vd *VectorData) implicateCycleHead(childNodesForNode map[string][]string,
		okNodes map[string]bool) error {
	involves := []string{}
	for name := range childNodesForNode {
		if !okNodes[name] && len(vd.mapItems[name].Referrers()) > 1 {
			involves = append(involves, "'" + name + "'")
		}
	}
	sort.Strings(involves)
	return fmt.Errorf("cycle detected in dependency graph; check references to %s",
		andList(involves))
}


func (vd *VectorData) CheckInStylesAndAttestations() error {
	if vd.styler == nil && vd.attester == nil {
		return fmt.Errorf("styles and attestations have not been configured")
	}
	err := vd.styler.checkConfiguration()
	if err != nil {
		return err
	}
	err = vd.attester.checkConfiguration()
	if err != nil {
		return err
	}
	for _, node := range vd.mapItems {
		style, attestation := node.styleAndAttestation()
		if style != nil || attestation != nil {
			if attestation == nil {
				err = vd.styler.resolveStyle(style)
			} else {
				err = vd.attester.resolveStyle(attestation, style)
			}
			if err != nil {
				break
			}
		}
	}
	return err
}

