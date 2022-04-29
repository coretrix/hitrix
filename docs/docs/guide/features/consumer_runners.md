# Consumer runners

Consumer runners enable you to quickly spin up BeeORM queue consumers easily.
There are 2 types of consumer.
- scalable
- non scalable


### ConsumerRunner - non scalable
Use ```queue.NewConsumerRunner(ctx)``` to make consumers which are not required to be able
to scale. 

This consumer works with following 4 interfaces:
- ConsumerOne (consumes items one by one)
- ConsumerMany (consumes items in batches)
- ConsumerOneByModulo (consumes items one by one using modulo)
- ConsumerManyByModulo (consumes items in batches using modulo)

### ScalableConsumerRunner - scalable
Use ```queue.NewScalableConsumerRunner(ctx, persistent redis pool)``` to make consumers which are required to be able
to scale.

This consumer works with following 2 interfaces:
- ConsumerOne (consumes items one by one)
- ConsumerMany (consumes items in batches)
