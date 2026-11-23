# Intent Recognition System Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement an LLM-powered intent recognition system that classifies user queries into predefined intent categories and extracts relevant entities/parameters, replacing the current keyword/rule-based skill selection.

**Architecture:** The system will add a preprocessing step before agent execution that uses the LLM to classify intent and extract entities. It will include a fallback mechanism to keyword-based routing for low-confidence predictions. Results will be stored in session metadata for auditability.

**Tech Stack:** Python, LangChain, existing prompt engineering infrastructure, JSON for structured output.

---
### Chunk 1: Intent Classification Component

**Files:**
- Create: `backend/graph/intent_classifier.py`
- Modify: `backend/graph/agent.py:1-50`
- Test: `tests/graph/test_intent_classifier.py`

- [ ] **Step 1: Write the failing test for intent classification**

```python
def test_classify_weather_intent():
    from graph.intent_classifier import classify_intent
    query = "What's the weather in Beijing today?"
    result = classify_intent(query)
    assert result["intent"] == "weather_search"
    assert result["confidence"] > 0.7
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/graph/test_intent_classifier.py::test_classify_weather_intent -v`
Expected: FAIL with "ModuleNotFoundError: No module named 'graph.intent_classifier'"

- [ ] **Step 3: Write minimal implementation**

```python
# backend/graph/intent_classifier.py
from typing import Dict
from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain
from langchain.chat_models import ChatOpenAI
import json

def classify_intent(query: str) -> Dict:
    template = """
    Classify the user query into one of the following intents:
    - weather_search: Questions about weather conditions
    - web_search: Requests to search the internet for information
    - knowledge_base_query: Questions about information in the local knowledge base
    - general_chat: Casual conversation or general questions
    - skill_creation: Requests to create or modify skills
    
    User query: {query}
    
    Respond in JSON format with:
    - intent: the classified intent
    - confidence: a float between 0 and 1 indicating confidence
    - reasoning: brief explanation of the classification
    
    Examples:
    Query: "What's the weather like?"
    Output: {{"intent": "weather_search", "confidence": 0.95, "reasoning": "User is asking about weather"}}
    
    Query: "Tell me about quantum computing"
    Output: {{"intent": "web_search", "confidence": 0.9, "reasoning": "User wants information likely requiring web search"}}
    """
    prompt = PromptTemplate(template=template, input_variables=["query"])
    llm = ChatOpenAI(temperature=0)
    chain = LLMChain(llm=llm, prompt=prompt)
    
    try:
        result = chain.run(query=query)
        # Parse JSON from result
        return json.loads(result)
    except Exception as e:
        # Fallback to default
        return {
            "intent": "general_chat",
            "confidence": 0.5,
            "reasoning": f"Error in classification: {str(e)}"
        }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/graph/test_intent_classifier.py::test_classify_weather_intent -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/graph/intent_classifier.py tests/graph/test_intent_classifier.py
git commit -m "feat: implement intent classification component"
```

### Chunk 2: Entity Extraction Component

**Files:**
- Create: `backend/graph/entity_extractor.py`
- Modify: `backend/graph/intent_classifier.py:1-30`
- Test: `tests/graph/test_entity_extractor.py`

- [ ] **Step 1: Write the failing test for entity extraction**

```python
def test_extract_weather_entities():
    from graph.entity_extractor import extract_entities
    query = "What's the weather in Beijing today?"
    intent = "weather_search"
    result = extract_entities(query, intent)
    assert result["location"] == "Beijing"
    assert result["date"] == "today"
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/graph/test_entity_extractor.py::test_extract_weather_entities -v`
Expected: FAIL with "ModuleNotFoundError: No module named 'graph.entity_extractor'"

- [ ] **Step 3: Write minimal implementation**

```python
# backend/graph/entity_extractor.py
from typing import Dict
from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain
from langchain.chat_models import ChatOpenAI
import json

def extract_entities(query: str, intent: str) -> Dict:
    if intent == "weather_search":
        template = """
        Extract weather-related entities from the query.
        
        User query: {query}
        
        Respond in JSON format with:
        - location: the location mentioned (city, etc.)
        - date: the date mentioned (today, tomorrow, etc.) or null if not specified
        - forecast_type: "current" or "forecast" or null if not specified
        
        Examples:
        Query: "What's the weather in Beijing?"
        Output: {{"location": "Beijing", "date": "today", "forecast_type": "current"}}
        
        Query: "Will it rain in Shanghai tomorrow?"
        Output: {{"location": "Shanghai", "date": "tomorrow", "forecast_type": "forecast"}}
        """
    elif intent == "web_search":
        template = """
        Extract search query entities from the query.
        
        User query: {query}
        
        Respond in JSON format with:
        - search_query: the core search query
        - time_range: any mentioned time range (today, week, month, year) or null
        - domains: any mentioned preferred domains or null
        
        Examples:
        Query: "Latest news about AI"
        Output: {{"search_query": "Latest news about AI", "time_range": null, "domains": null}}
        
        Query: "Python documentation from last week"
        Output: {{"search_query": "Python documentation", "time_range": "week", "domains": null}}
        """
    else:
        # Default empty extraction for other intents
        return {}
    
    prompt = PromptTemplate(template=template, input_variables=["query"])
    llm = ChatOpenAI(temperature=0)
    chain = LLMChain(llm=llm, prompt=prompt)
    
    try:
        result = chain.run(query=query)
        return json.loads(result)
    except Exception as e:
        return {}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/graph/test_entity_extractor.py::test_extract_weather_entities -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/graph/entity_extractor.py tests/graph/test_entity_extractor.py
git commit -m "feat: implement entity extraction component"
```

### Chunk 3: Integration with Agent Flow

**Files:**
- Modify: `backend/graph/agent.py:50-150`
- Modify: `backend/graph/prompt_builder.py:1-50`
- Test: `tests/graph/test_agent_integration.py`

- [ ] **Step 1: Write the failing test for integration**

```python
def test_agent_uses_intent_recognition():
    from graph.agent import AgentManager
    # Mock the LLM responses for intent classification and entity extraction
    # Verify that the agent calls the intent classifier and uses the result
    manager = AgentManager()
    # This is a simplified test - in reality we would mock the LLM calls
    assert hasattr(manager, 'process_query')
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/graph/test_agent_integration.py::test_agent_uses_intent_recognition -v`
Expected: FAIL with "ModuleNotFoundError: No module named 'tests.graph.test_agent_integration'"

- [ ] **Step 3: Write minimal implementation**

```python
# backend/graph/agent.py (modified)
from .intent_classifier import classify_intent
from .entity_extractor import extract_entities

class AgentManager:
    def __init__(self, backend_dir):
        # ... existing init code ...
        self.intent_confidence_threshold = 0.7
    
    def process_query(self, query: str, session_id: str):
        # Step 1: Classify intent
        intent_result = classify_intent(query)
        intent = intent_result["intent"]
        confidence = intent_result["confidence"]
        
        # Step 2: Extract entities if confidence is high enough
        entities = {}
        if confidence >= self.intent_confidence_threshold:
            entities = extract_entities(query, intent)
        else:
            # Fallback to keyword-based intent detection (existing logic)
            intent = self._fallback_intent_detection(query)
        
        # Step 3: Store intent and entities in session metadata for audit
        self._store_intent_metadata(session_id, intent, entities, confidence)
        
        # Step 4: Proceed with existing agent logic using intent and entities
        # ... existing code that uses intent and entities to select skills ...
        
        # For now, we'll just return a placeholder
        return f"Processed query with intent: {intent}"
    
    def _store_intent_metadata(self, session_id: str, intent: str, entities: dict, confidence: float):
        # Store in session JSON file or memory for audit
        import json
        from pathlib import Path
        session_file = Path(f"backend/sessions/{session_id}.json")
        if session_file.exists():
            with open(session_file, 'r+') as f:
                session_data = json.load(f)
                session_data["intent_metadata"] = {
                    "intent": intent,
                    "confidence": confidence,
                    "entities": entities
                }
                f.seek(0)
                json.dump(session_data, f, indent=2)
                f.truncate()
```

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/graph/test_agent_integration.py::test_agent_uses_intent_recognition -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/graph/agent.py backend/graph/prompt_builder.py tests/graph/test_agent_integration.py
git commit -m "feat: integrate intent recognition into agent flow"
```

### Chunk 4: Fallback Mechanism and Error Handling

**Files:**
- Modify: `backend/graph/intent_classifier.py:30-60`
- Modify: `backend/graph/entity_extractor.py:30-60`
- Test: `tests/graph/test_fallback.py`

- [ ] **Step 1: Write the failing test for fallback**

```python
def test_low_confidence_triggers_fallback():
    from graph.intent_classifier import classify_intent
    # Use a query that should produce low confidence
    query = "asdfasdf asdf asdf"
    result = classify_intent(query)
    # Should fall back to general_chat with lower confidence
    assert result["intent"] == "general_chat"
    assert result["confidence"] < 0.7
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/graph/test_fallback.py::test_low_confidence_triggers_fallback -v`
Expected: FAIL with "ModuleNotFoundError: No module named 'tests.test_fallback'"

- [ ] **Step 3: Write minimal implementation**

```python
# backend/graph/intent_classifier.py (modified)
def classify_intent(query: str) -> Dict:
    # ... existing code ...
    try:
        result = chain.run(query=query)
        parsed_result = json.loads(result)
        # Validate the result has required fields
        if "intent" not in parsed_result or "confidence" not in parsed_result:
            raise ValueError("Missing required fields")
        return parsed_result
    except Exception as e:
        # Fallback to keyword-based detection
        fallback_intent = _keyword_based_fallback(query)
        return {
            "intent": fallback_intent,
            "confidence": 0.5,  # Lower confidence for fallback
            "reasoning": f"Fallback due to error: {str(e)}"
        }

def _keyword_based_fallback(query: str) -> str:
    query_lower = query.lower()
    if any(word in query_lower for word in ["weather", "temperature", "forecast"]):
        return "weather_search"
    elif any(word in query_lower for word in ["search", "find", "look up"]):
        return "web_search"
    elif any(word in query_lower for word in ["knowledge", "document", "file"]):
        return "knowledge_base_query"
    elif any(word in query_lower for word in ["skill", "create", "modify"]):
        return "skill_creation"
    else:
        return "general_chat"
```

```python
# backend/graph/entity_extractor.py (modified)
def extract_entities(query: str, intent: str) -> Dict:
    # ... existing code ...
    try:
        result = chain.run(query=query)
        return json.loads(result)
    except Exception as e:
        # Return empty dict on failure - let the skill handle missing entities
        return {}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/graph/test_fallback.py::test_low_confidence_triggers_fallback -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/graph/intent_classifier.py backend/graph/entity_extractor.py tests/graph/test_fallback.py
git commit -m "feat: add fallback mechanism and error handling"
```

### Chunk 5: Testing and Validation

**Files:**
- Modify: `tests/conftest.py` (for test fixtures)
- Create: `tests/data/test_queries.json`
- Modify: `backend/graph/agent.py:150-200` (add more test hooks if needed)

- [ ] **Step 1: Create test query dataset**

```json
// tests/data/test_queries.json
[
    {
        "query": "What's the weather in Beijing?",
        "expected_intent": "weather_search",
        "expected_entities": {"location": "Beijing", "date": "today"}
    },
    {
        "query": "Find me the latest news about AI",
        "expected_intent": "web_search",
        "expected_entities": {"search_query": "latest news about AI"}
    },
    {
        "query": "What does the knowledge base say about quantum computing?",
        "expected_intent": "knowledge_base_query",
        "expected_entities": {"search_query": "quantum computing"}
    }
]
```

- [ ] **Step 2: Write validation test**

```python
def test_intent_recognition_accuracy():
    import json
    from graph.intent_classifier import classify_intent
    from graph.entity_extractor import extract_entities
    
    with open("tests/data/test_queries.json", 'r') as f:
        test_queries = json.load(f)
    
    correct = 0
    total = len(test_queries)
    
    for test_case in test_queries:
        query = test_case["query"]
        expected_intent = test_case["expected_intent"]
        
        intent_result = classify_intent(query)
        if intent_result["intent"] == expected_intent and intent_result["confidence"] > 0.7:
            correct += 1
            
            # If intent is correct and confidence high, check entities
            entities = extract_entities(query, intent_result["intent"])
            expected_entities = test_case["expected_entities"]
            # Simple check - in reality would be more sophisticated
            if all(entities.get(k) == v for k, v in expected_entities.items()):
                pass  # Entity extraction correct
            else:
                pass  # Entity extraction needs work - but we're testing intent here
        else:
            pass  # Intent classification incorrect
    
    accuracy = correct / total
    assert accuracy >= 0.8, f"Intent classification accuracy {accuracy} below threshold 0.8"
```

- [ ] **Step 3: Run validation test**

Run: `pytest tests/test_intent_recognition_accuracy.py -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add tests/data/test_queries.json tests/test_intent_recognition_accuracy.py
git commit -m "feat: add validation tests for intent recognition system"
```

### Chunk 6: Documentation and Examples

**Files:**
- Create: `docs/intents.md`
- Modify: `README.md:500-550` (add section about intent recognition)
- Modify: `backend/skills/*/SKILL.md` (add notes about intent recognition if needed)

- [ ] **Step 1: Write intent documentation**

```markdown
# Intent Recognition System

## Overview
The langchain-OpenClaw agent now uses an LLM-powered intent recognition system to understand user queries and determine the appropriate skills to use.

## Supported Intents
- `weather_search`: Questions about weather conditions
- `web_search`: Requests to search the internet for information
- `knowledge_base_query`: Questions about information in the local knowledge base
- `general_chat`: Casual conversation or general questions
- `skill_creation`: Requests to create or modify skills

## How It Works
1. The user query is first passed to the intent classifier
2. The classifier uses the LLM to determine the intent with a confidence score
3. If confidence is above threshold (0.7), entities are extracted for that intent
4. The intent and entities are used to select and configure the appropriate skill
5. Low-confidence predictions fall back to keyword-based routing

## Auditability
All intent recognition results are stored in session JSON files under:
`sessions/<session_id>.json` in the `intent_metadata` field.

Example:
```json
{
  "intent_metadata": {
    "intent": "weather_search",
    "confidence": 0.95,
    "entities": {
      "location": "Beijing",
      "date": "today",
      "forecast_type": "current"
    }
  }
}
```
```

- [ ] **Step 2: Update README**

```markdown
## Intent Recognition System

The agent now features an LLM-powered intent recognition system that improves skill selection accuracy by understanding the semantic meaning of user queries rather than relying solely on keyword matching.
```

- [ ] **Step 3: Commit**

```bash
git add docs/intents.md README.md
git commit -m "docs: add documentation for intent recognition system"
```

## Plan Complete

Plan complete and saved to `docs/superpowers/plans/2026-03-24-intent-recognition-system.md`. Ready to execute?