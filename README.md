# Rule Engine
Design a Rule Engine to support a dynamic and configurable rule-based system. The Rule Engine should allow the configuration of different types of rules, independent execution of rules, and storage of rule evaluation results. These rules should be possible to write as code.

Consider the following requirements:

- Rule Types: The system should support multiple types of rules, such as validation rules, scoring rules, or conditional rules.
- Independent Execution: Each rule should be capable of independent execution, ensuring that the failure or success of one rule does not impact others.
- Result Storage: Results of rule evaluations should be stored for later analysis.
- Flexibility: The Rule Engine should be flexible enough to accommodate new rule types without significant changes to the core system.

Your task is to design and implement a Rule Engine that meets these requirements. Please provide:

- A high-level design of the Rule Engine, including the key components and their interactions.
- An example implementation of a rule type (e.g., validation rule) with configurable parameters.
- Demonstration of how the Rule Engine can execute multiple rules independently.
- Considerations for result storage and retrieval.

Consider following requirement for design extension

- Having multiple rule groups, each with x number of rules. You can define whether the rules inside a rule group should be executed sequentially or concurrently. You can define whether to exit the rule engine after getting a result of a rule group - like skip further rules or rule groups.

Include unit tests to validate the functionality of the Rule Engine.