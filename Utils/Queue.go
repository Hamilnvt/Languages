package Utils

import (
  "fmt"
  "errors"
)

type Queue[S any] struct {
  queue []S
}

func (queue Queue[S]) String() (res string) {
  for i := 0; i < len(queue.queue); i++ {
    res += fmt.Sprintf("%v. %v\n", i, queue.queue[i])
  }
  return
}

func (queue Queue[S]) First() (S, error) {
  if queue.IsEmpty() {
    var empty S
    return empty, errors.New("Empty queue")
  } else {
    return queue.queue[0], nil
  }
}

func (queue Queue[S]) IsEmpty() bool {
  return len(queue.queue) == 0
}

func (queue *Queue[S]) Enqueue(elt S) {
  queue.queue = append(queue.queue, elt)
}

func (queue *Queue[S]) Dequeue() (S, error) {
  if queue.IsEmpty() {
    var empty S
    return empty, errors.New("ERROR: Cannot dequeue from empty queue")
  } else {
    elt, _ := queue.First()
    queue.queue = queue.queue[1:]
    return elt, nil
  }
}
