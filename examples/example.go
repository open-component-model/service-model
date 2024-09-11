package examples

import _ "embed"

//go:embed  descriptors/MSPGardener.yaml
var MSPGardener string

//go:embed  descriptors/MSPHana.yaml
var MSPHana string

//go:embed  descriptors/MSPSteampunk.yaml
var MSPSteampunk string

//go:embed  descriptors/SvcAbap.yaml
var SvcAbap string

//go:embed  descriptors/SvcGardenCluster22.yaml
var SvcGardenCluster22 string

//go:embed  descriptors/SvcGardenCluster23.yaml
var SvcGardenCluster23 string

//go:embed  descriptors/ContractK8sCluster23.yaml
var ContractK8sCluster23 string

//go:embed  descriptors/ContractK8sCluster22.yaml
var ContractK8sCluster22 string

//go:embed  descriptors/SvcHana.yaml
var SvcHana string
