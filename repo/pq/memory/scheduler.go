package memory

/*
import (
	"context"
	"fmt"
	"time"
)

func (a *Adapter) delayedScheduler(exit chan struct{}) {
	defer a.wg.Done()

	secTicker := time.NewTicker(time.Second)
	defer secTicker.Stop()

	recapTicker := time.NewTicker(time.Hour)
	defer recapTicker.Stop()

loop:
	for {
		select {
		case <-exit:
			break loop

		case t := <-secTicker.C:
			ts := t.Unix()
			for {
				// if a.pqDelayed.Len() == 0 {
				// 	break
				// }

				jitem := a.pqDelayed.Peek()
				if jitem == nil || jitem.heapPriority > ts {
					break
				}

				jitems, _ := a.pqDelayed.Pop(1)
				if len(jitems) == 1 {
					jitem = jitems[0]
					_, err := a.Push(context.Background(), jitem.Topic, jitem.Priority, jitem.JobID, time.Time{})
					if err != nil {
						fmt.Print(err)
					}
				}
			}

		case <-recapTicker.C:
			a.pqDelayed.Recap(100, 1.3)
			for _, pq := range a.pqByTopic {
				pq.Recap(100, 1.3)
			}
		}
	}
}
*/
