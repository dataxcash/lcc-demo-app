# Multi-Instance Backend API

## Overview

The demo app backend now supports registering and managing **multiple SDK instances** simultaneously. Each instance gets its own unique ID, allowing you to:

- Register multiple products in one go
- Register the same product with different versions as separate instances
- Manage each instance independently
- Test and delete instances individually

## Implementation Details

### Data Structures

```go
// Instance represents a registered SDK instance
type Instance struct {
    InstanceID   string    `json:"instance_id"`    // Unique instance ID from LCC
    ProductID    string    `json:"product_id"`     // Product being tested
    Version      string    `json:"version"`        // Product version
    RegisteredAt time.Time `json:"registered_at"`  // Registration timestamp
}

// Server maintains:
// - instances: map[string]*Instance          // All registered instances
// - instanceKeys: map[string]*auth.KeyPair   // RSA keys per instance
// - clients: map[string]*lccclient.Client    // SDK clients (backward compat)
```

### Storage Strategy

**Unique Key Format**: `{product_id}:{version}:{instance_id}`

This allows:
- Multiple versions of same product: `product-x:1.0.0:uuid1`, `product-x:2.0.0:uuid2`
- Instances with same product ID but different versions tracked separately
- Full flexibility for version-based testing

### API Endpoints

#### 1. Register Instance
```
POST /api/instance/register
```

**Request**:
```json
{
    "product_id": "string",
    "version": "string (optional, defaults to 1.0.0)",
    "lcc_url": "string (optional, uses saved config if not provided)"
}
```

**Response**:
```json
{
    "success": true,
    "instance_id": "uuid-string",
    "product_id": "string",
    "version": "string",
    "registered_at": "2025-11-22T13:00:00Z",
    "message": "Instance registered successfully"
}
```

**Behavior**:
- Creates new Instance struct
- Stores in `instances` map with unique key
- Stores RSA KeyPair in `instanceKeys` map
- Maintains backward compatibility with `clients` map

#### 2. List All Instances
```
GET /api/instances
```

**Response**:
```json
{
    "instances": [
        {
            "instance_id": "uuid-1",
            "product_id": "product-x",
            "version": "1.0.0",
            "registered_at": "2025-11-22T12:50:00Z"
        },
        {
            "instance_id": "uuid-2",
            "product_id": "product-x",
            "version": "2.0.0",
            "registered_at": "2025-11-22T13:00:00Z"
        },
        {
            "instance_id": "uuid-3",
            "product_id": "product-y",
            "version": "1.0.0",
            "registered_at": "2025-11-22T13:05:00Z"
        }
    ]
}
```

**Features**:
- Returns all registered instances
- RW mutex protection for thread-safe reads
- No filter needed - returns everything

#### 3. Clear Instance (Enhanced)
```
POST /api/instance/clear
```

**Request (Clear specific instance)**:
```json
{
    "product_id": "string",
    "instance_id": "string (optional)"
}
```

**Behavior**:
- If `instance_id` provided: Delete only that specific instance
- If `instance_id` omitted: Delete all instances for the product
- Updates both `instances` and `instanceKeys` maps
- Maintains backward compatibility with `clients` map

**Response**:
```json
{
    "success": true,
    "message": "Instance cleared successfully"
}
```

#### 4. Test Instance (Unchanged)
```
POST /api/instance/test
```

**Request**:
```json
{
    "product_id": "string",
    "feature_id": "string"
}
```

Works with any registered instance using product_id lookup.

#### 5. Get Instance Status (Existing)
```
GET /api/instance/status?product_id=string
```

Works with existing instances via product_id.

## Concurrency & Thread Safety

All operations use RW mutex (`sync.RWMutex`):

```go
s.mu.Lock()      // For writes (register, delete)
s.mu.Unlock()

s.mu.RLock()     // For reads (list)
s.mu.RUnlock()
```

**Operations**:
- **Register**: Lock, insert into maps, unlock
- **List**: RLock, copy instances, RUnlock
- **Delete**: Lock, remove from maps, unlock

## Backward Compatibility

The implementation maintains full backward compatibility:

1. **Single Product Lookup**: Still works via old `clients` map
2. **Legacy REST Calls**: All existing endpoints unchanged
3. **Graceful Degradation**: If instance not found, returns sensible errors

## Example Workflows

### Workflow 1: Multi-Version Testing

```bash
# Register v1.0.0
curl -X POST http://localhost:9144/api/instance/register \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "my-product",
    "version": "1.0.0",
    "lcc_url": "http://localhost:7086"
  }'

# Response: Instance ID = abc-123

# Register v2.0.0
curl -X POST http://localhost:9144/api/instance/register \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "my-product",
    "version": "2.0.0",
    "lcc_url": "http://localhost:7086"
  }'

# Response: Instance ID = abc-456

# List all instances
curl http://localhost:9144/api/instances
# Returns both v1 and v2 instances

# Test v1.0.0
curl -X POST http://localhost:9144/api/instance/test \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "my-product",
    "feature_id": "feature-x"
  }'
# Uses v1.0.0 client

# Delete only v2.0.0
curl -X POST http://localhost:9144/api/instance/clear \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "my-product",
    "instance_id": "abc-456"
  }'
# v1.0.0 instance remains
```

### Workflow 2: Batch Registration

```bash
# Frontend selects: ProductA, ProductB, ProductC
# Frontend sets version: 1.5.0

# Backend makes 3 sequential calls:
POST /api/instance/register (ProductA, 1.5.0) → Instance ID: xxx-111
POST /api/instance/register (ProductB, 1.5.0) → Instance ID: xxx-222
POST /api/instance/register (ProductC, 1.5.0) → Instance ID: xxx-333

# Result in instances map:
{
  "product-a:1.5.0:xxx-111": {...},
  "product-b:1.5.0:xxx-222": {...},
  "product-c:1.5.0:xxx-333": {...}
}
```

### Workflow 3: Independent Quota Tracking

Each instance maintains independent quota:

```
ProductX v1.0.0 (abc-123):
  - Consume 100 units → 900 remaining
  
ProductX v1.0.0 (abc-456) [different instance]:
  - Consume 100 units → 900 remaining (independent)
```

## Limitations

**Current**:
- Instances stored in memory (lost on server restart)
- No persistent storage between sessions
- No export/import of configurations

**Future**:
- Database persistence (PostgreSQL/MySQL)
- Embedded DB (SQLite for demo)
- Export instance configs as JSON/YAML
- Batch import from file

## Troubleshooting

### Multiple registrations of same product show different instance IDs

**Expected behavior** - each registration creates a new instance with a new ID. This is intentional for multi-instance support.

### Deleted instance still appears in list

**Check**: Are you deleting with the correct `instance_id`? Without instanceID, it deletes all instances for that product.

### Instance test fails after adding multi-instance support

**Check**: You're testing with the right product_id? Multi-instance uses product_id for test routing, not instance_id directly (for backward compat).

## Implementation Code

### Key Changes in `server.go`:
```go
type Instance struct {
    InstanceID   string    `json:"instance_id"`
    ProductID    string    `json:"product_id"`
    Version      string    `json:"version"`
    RegisteredAt time.Time `json:"registered_at"`
}

type Server struct {
    // ...existing fields...
    instances    map[string]*Instance       // NEW: Multi-instance storage
    instanceKeys map[string]*auth.KeyPair   // NEW: Keys per instance
}
```

### Key Changes in `instance_handler.go`:
```go
// In handleInstanceRegister:
instanceKey := fmt.Sprintf("%s:%s:%s", req.ProductID, req.Version, instanceID)
s.instances[instanceKey] = &Instance{...}
s.instanceKeys[instanceKey] = kp

// In handleListInstances (NEW):
// Copy all instances from map and return as JSON array

// In handleInstanceClear (UPDATED):
// Support optional instance_id parameter
// Delete specific instance or all for product
```

## Performance Characteristics

- **Register**: O(1) - map insert
- **List**: O(n) - iterate all instances
- **Delete**: O(n) - search by instance_id (if specified)
- **Lookup**: O(1) - direct map access

Where n = number of instances (typically < 1000 for demo)
