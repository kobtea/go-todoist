package todoist

type RelationClient struct {
	*Client
}

type ItemRelations struct {
	//Users map[ID]User
	Projects map[ID]Project
	Labels   map[ID]Label
}

func (c RelationClient) Items(items []Item) ItemRelations {
	res := ItemRelations{Projects: map[ID]Project{}, Labels: map[ID]Label{}}
	for _, item := range items {
		if _, ok := res.Projects[item.ProjectID]; !ok {
			p := c.Project.Resolve(item.ProjectID)
			if p != nil {
				res.Projects[item.ProjectID] = *p
			}
		}
		for _, id := range item.Labels {
			if _, ok := res.Labels[id]; !ok {
				l := c.Label.Resolve(id)
				if l != nil {
					res.Labels[id] = *l
				}
			}
		}
	}
	return res
}
