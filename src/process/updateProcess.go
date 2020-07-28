package process

// UpdateProcInfo updates the fields for a process.
func (p *Process) UpdateProcInfo() {
	tempBackground, err := p.Proc.Background()
	if err == nil {
		p.Background = tempBackground
	}

	tempForeground, err := p.Proc.Foreground()
	if err == nil {
		p.Foreground = tempForeground
	}

	tempIsRunning, err := p.Proc.IsRunning()
	if err == nil {
		p.IsRunning = tempIsRunning
	}

	tempCPUPercent, err := p.Proc.CPUPercent()
	if err == nil {
		p.CPUPercent = tempCPUPercent
	}

	tempChildren, err := p.Proc.Children()
	if err == nil {
		p.Children = tempChildren
	}

	tempCreateTime, err := p.Proc.CreateTime()
	if err == nil {
		p.CreateTime = tempCreateTime
	}

	tempGids, err := p.Proc.Gids()
	if err == nil {
		p.Gids = tempGids
	}

	tempMemInfo, err := p.Proc.MemoryInfo()
	if err == nil {
		p.MemoryInfo = tempMemInfo
	}

	tempMemPerc, err := p.Proc.MemoryPercent()
	if err == nil {
		p.MemoryPercent = tempMemPerc
	}

	tempName, err := p.Proc.Name()
	if err == nil {
		p.Name = tempName
	}

	tempNice, err := p.Proc.Nice()
	if err == nil {
		p.Nice = tempNice
	}

	tempCtx, err := p.Proc.NumCtxSwitches()
	if err == nil {
		p.NumCtxSwitches = tempCtx
	}

	tempNumThreads, err := p.Proc.NumThreads()
	if err == nil {
		p.NumThreads = tempNumThreads
	}

	tempPageFault, err := p.Proc.PageFaults()
	if err == nil {
		p.PageFault = tempPageFault
	}

	tempStatus, err := p.Proc.Status()
	if err == nil {
		p.Status = tempStatus
	}

	tempExe, err := p.Proc.Exe()
	if err == nil {
		p.Exe = tempExe
	} else {
		p.Exe = "NA"
	}

	tempAffinity, err := p.Proc.CPUAffinity()
	if err == nil {
		p.CPUAffinity = tempAffinity
	}
}

func (p *Process) UpdateProcForVisual() {
	tempForeground, err := p.Proc.Foreground()
	if err == nil {
		p.Foreground = tempForeground
	}

	tempCPUPercent, err := p.Proc.CPUPercent()
	if err == nil {
		p.CPUPercent = tempCPUPercent
	}

	tempMemPerc, err := p.Proc.MemoryPercent()
	if err == nil {
		p.MemoryPercent = tempMemPerc
	}

	tempNumThreads, err := p.Proc.NumThreads()
	if err == nil {
		p.NumThreads = tempNumThreads
	}

	tempCreateTime, err := p.Proc.CreateTime()
	if err == nil {
		p.CreateTime = tempCreateTime
	}

	tempStatus, err := p.Proc.Status()
	if err == nil {
		p.Status = tempStatus
	}
}
