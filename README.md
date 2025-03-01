# saga-pattern

This is an implementation of the Saga Pattern in a single microservice, which is used to manage the checkout workflow.

This implementation includes 3 steps: Payment, Inventory, Shipping. Each of which implements the Step interface with the Do (to perform an action) and Compensate (to roll back to the previous state in case of an error) methods.

The Saga structure itself is implemented as a set of steps that are cyclically called, and in case of an error on any of them, Compensate is called for all previous ones. It has methods Execute (for executing the workflow) and Rollback (for the process of calling Compensate for the previous steps).