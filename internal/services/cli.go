package services

import "time"

func (d *GokitService) CreateCompetition() error {
	time.Sleep(7 * time.Second)

	d.logger.Info("Create Competition Successfully")

	return nil
}

func (d *GokitService) CreateTeam() error {
	time.Sleep(7 * time.Second)
	d.logger.Info("Create Team Successfully")
	return nil
}
