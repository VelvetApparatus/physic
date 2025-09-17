## domain layer: aggregate

#### 1. Validate only business errors. Business errors determinate business flow.
#### 2. All validations should be encapsulated in `rules`. Rules make business logic more `transparent`.
#### 3. If Aggregate has value objects in itself: it should use validation methods in value_objects package.
#### 4. Aggregates can have infrastructure packages in itself: but usage of them should be hidden in DI. Best practise: use custom context and provide all needs in them.