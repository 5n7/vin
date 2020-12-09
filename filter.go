package vin

func (v Vin) Filter(filter func(app App) bool) *Vin {
	apps := make([]App, 0)
	for _, app := range v.Apps {
		if filter(app) {
			apps = append(apps, app)
		}
	}
	v.Apps = apps
	return &v
}

func (v *Vin) FilterByHost(host string) *Vin {
	return v.Filter(func(app App) bool {
		if len(app.Hosts) == 0 {
			return true
		}

		for _, h := range app.Hosts {
			if h == host {
				return true
			}
		}
		return false
	})
}

func (v *Vin) FilterByPriority(minPriority int) *Vin {
	return v.Filter(func(app App) bool {
		return app.Priority >= minPriority
	})
}

func (v *Vin) FilterByRepo(repos []string) *Vin {
	return v.Filter(func(app App) bool {
		for _, repo := range repos {
			if app.Repo == repo {
				return true
			}
		}
		return false
	})
}
