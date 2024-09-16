package typehandler

import (
	"github.com/mandelsoft/goutils/set"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/cmds/ocm/common/data"
	"ocm.software/ocm/cmds/ocm/common/processing"
)

type NormObject interface {
	common.HistorySource
	GetKey() common.NameVersion
	CreateContinue() NormObject
}

type NormObjects = Objects[NormObject]

var Normalize = processing.Transform(NormalizeFunction)

func NormalizeFunction(s data.Iterable) data.Iterable {
	objs := ObjectSlice[NormObject](s)

	uses := set.Set[common.NameVersion]{}
	used := set.Set[common.NameVersion]{}
	depth := map[common.NameVersion]int{}
	found := set.Set[common.NameVersion]{}

	var skip common.History
	var base NormObject

	for _, o := range objs {
		key := o.GetKey()
		hist := o.GetHistory()
		if base != nil && hist.HasPrefix(skip) {
			uses.Add(base.GetKey())
		}
		if len(hist) > 0 {
			used.Add(key)
		}
		base = o
		skip = o.GetHistory().Append(key)
	}

	skip = nil
	for _, o := range objs {
		key := o.GetKey()
		hist := o.GetHistory()
		if d, ok := depth[key]; (!ok || (d > len(hist))) && !(len(hist) == 0 && used.Has(key)) {
			depth[key] = len(hist)
		}
	}

	var result NormObjects
	skip = nil
	for _, o := range objs {
		key := o.GetKey()
		hist := o.GetHistory()
		if skip != nil && hist.HasPrefix(skip) {
			continue
		}
		skip = nil
		_, ok := found[key]
		if ok || len(hist) != depth[key] || found.Has(key) {
			skip = hist.Append(key)
			base = o
			if len(hist) == 0 {
				continue
			}
		} else {
			found.Add(key)
		}
		result = append(result, o)
		if skip != nil && uses.Has(key) {
			result = append(result, o.CreateContinue())
		}
	}
	return result
}
