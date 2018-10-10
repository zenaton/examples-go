package workflow

// VersionDefinition represents a versioned workflow definition
type VersionDefinition struct {
	name     string
	versions []*Definition
}

// Version will return a VersionDefinition.
// Trying to update a workflow implementation while some instances are running will trigger a panic.
// This is because Zenaton has currently no way to know how to use an updated implementation to continue execution of
// already running workflow instances.
// Zenaton lets you easily update a workflow as described below, but keep in mind your modifications will apply only to
// new instances.
//
// Let's assume you currently have a workflow called MyWorkflow. As you will have different versions running, you MUST version your
// sources. This is how to do it:
// 		rename your existing workflow MyWorkflow to MyWorkflow_v0
// 		create your new version MyWorkflow_v1
// 		create a new MyWorkflow workflow with the provided Version function
//
// For Example:
//		import "github.com/zenaton/zenaton-go/v1/zenaton/workflow"
//		var MyWorkflow = workflow.Version("MyWorkflow", []*workflow.Definition{
//			MyWorkflow_v0,
//			MyWorkflow_v1,
//		})
// After that, if you write another version MyWorkflow_v2, you just have to add it to your sources and add it to MyWorkflow.
//
// You do NOT have to change the implementation of your client that will still use MyWorkflow class (eg. to launch a
// workflow, you will still use MyWorkflow.New(...).Dispatch() ).
//
// The _v convention is only a suggestion. You can decide on a different naming convention.
func Version(name string, workflowDefinitions []*Definition) *VersionDefinition {

	if len(workflowDefinitions) == 0 {
		panic("must provide at least one workflow definition to Version")
	}

	versionedDef := &VersionDefinition{
		name:     name,
		versions: workflowDefinitions,
	}

	UnsafeManager.setVersionDef(name, versionedDef)

	return versionedDef
}

// NewInstance will return an instance of the last added workflow Definition.
// Any arguments you pass in here will be given to your Init() function. Thus, the arguments must be of the same
// number and type as your Init function expects.
func (vd *VersionDefinition) NewInstance(args ...interface{}) *Instance {

	var instance *Instance
	if len(args) > 0 {
		instance = vd.versions[len(vd.versions)-1].New(args...)
	} else {
		instance = vd.versions[len(vd.versions)-1].New()
	}

	instance.canonical = vd.name
	return instance

}

func (vd *VersionDefinition) getCurrentDefinition() *Definition {
	return vd.versions[len(vd.versions)-1]
}

func (vd *VersionDefinition) getInitialDefinition() *Definition {
	return vd.versions[0]
}

// WhereID takes an id (of a workflow instance) and returns a QueryBuilder.
// The QueryBuilder allows you to Find, Kill, Pause, and Resume workflow instances by id. You can also Send an event
// to a workflow.
func (vd *VersionDefinition) WhereID(id string) *QueryBuilder {
	return newBuilder(vd.name).whereID(id)
}
