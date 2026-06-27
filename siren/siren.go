package siren

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

type Siren struct {
	wavBytes   []byte
	wavSpec    sdl.AudioSpec
	volume     float32
	soundLoops uint
}

func New(path string, volume float32) (*Siren, error) {
	var wavSpec sdl.AudioSpec
	wav, err := sdl.LoadWAV(path, &wavSpec)
	if err != nil {
		return nil, err
	}

	return &Siren{
		wavBytes:   wav,
		wavSpec:    wavSpec,
		volume:     volume,
		soundLoops: 24,
	}, nil
}

func (siren *Siren) Start() (chan error, error) {
	devices, err := sdl.GetAudioPlaybackDevices()
	if err != nil {
		return nil, err
	}

	audioStream := devices[0].OpenAudioDeviceStream(&siren.wavSpec, 0)
	err = audioStream.ResumeDevice()
	if err != nil {
		audioStream.Destroy()
		return nil, err
	}
	audioStream.SetGain(siren.volume)

	for range siren.soundLoops {
		err = audioStream.PutData(siren.wavBytes)
		if err != nil {
			audioStream.Destroy()
			return nil, err
		}
	}
	err = audioStream.Flush()
	if err != nil {
		audioStream.Destroy()
		return nil, err
	}

	resultChan := make(chan error)
	go func(rc chan error) {
		defer audioStream.Destroy()
		for {
			av, err := audioStream.Available()
			if err != nil {
				rc <- err
			}
			if av > 0 {
				sdl.Delay(100)
			} else {
				break
			}
		}
		rc <- nil
		close(rc)
	}(resultChan)

	return nil, nil
}
