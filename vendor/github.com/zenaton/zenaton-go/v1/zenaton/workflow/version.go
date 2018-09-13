package workflow

import "github.com/zenaton/zenaton-go/v1/zenaton/interfaces"

//todo:
type VersionDefinition struct {
	name     string
	versions []*Definition
}

func Version(name string, workflowDefinitions []*Definition) *VersionDefinition {

	if len(workflowDefinitions) == 0 {
		panic("must provide at least one workflow definition to Version")
	}

	versionedDef := &VersionDefinition{
		name:     name,
		versions: workflowDefinitions,
	}

	Manager.setVersionDef(name, versionedDef)

	return versionedDef
}

func (vd *VersionDefinition) NewInstance(handlers ...interfaces.Handler) *Instance {

	if len(handlers) > 1 {
		panic("must pass at maximum one handler to VersionDefinition.NewInstance()")
	}

	var instance *Instance
	if len(handlers) == 1 {
		h := handlers[0]
		validateHandler(h)
		instance = vd.versions[len(vd.versions)-1].NewInstance(h)
	} else {
		instance = vd.versions[len(vd.versions)-1].NewInstance()
	}
	instance.setCanonical(vd.name)
	return instance

}

func (vd *VersionDefinition) getCurrentDefinition() *Definition {
	return vd.versions[len(vd.versions)-1]
}

func (vd *VersionDefinition) getInitialDefinition() *Definition {
	return vd.versions[0]
}

//todo: this needs testing
func (vd *VersionDefinition) WhereID(id string) *QueryBuilder {
	return NewBuilder(vd.name).WhereID(id)
}
