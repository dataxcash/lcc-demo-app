# Multi-Instance SDK Setup

## Overview

The enhanced Setup page now supports registering **multiple SDK instances** simultaneously, removing the previous limitation of one instance per product.

## Key Improvements

### 1. **Batch Registration**
- Register multiple products in a single operation
- Each product gets its own unique instance ID
- Progress tracking for each registration

### 2. **Version Testing**
- Register the same product with different versions
- Example: `product-x v1.0.0` and `product-x v2.0.0` as separate instances
- Each version maintains independent quota tracking

### 3. **Instance Dashboard**
- View all registered instances in one place
- Shows instance ID, product ID, version, and status
- Real-time status updates

### 4. **Multi-Instance Management**
- Select and test individual instances
- Delete instances without affecting others
- Test connections independently

## Use Cases

### A/B Testing
```
Register two versions of the same product:
- Product A (v1.0.0) → Instance ID: xxx-123
- Product A (v1.1.0) → Instance ID: xxx-456

Compare behavior and quota consumption between versions
```

### Load Testing
```
Register multiple instances of the same product:
- Product B (v1.0.0) → Instance ID: yyy-111
- Product B (v1.0.0) → Instance ID: yyy-222
- Product B (v1.0.0) → Instance ID: yyy-333

Run concurrent load tests with 3 independent quota buckets
```

### Multi-Product Testing
```
Register all products at once:
- Product A (v1.0.0) → Instance ID: aaa-111
- Product B (v2.0.0) → Instance ID: bbb-222
- Product C (v1.5.0) → Instance ID: ccc-333

Test multiple products in parallel using a single registration UI
```

## UI Features

### Step 1: Register Products
1. Enter LCC Server URL (cached from previous configuration)
2. Select multiple products using checkboxes
3. Set version (applies to all selected products)
4. Click "Register Selected Products"
5. Monitor progress with progress bars

### Step 2: View Dashboard
- See all registered instances
- Instance ID, Product ID, Version, Status displayed
- Auto-updates after each registration

### Step 3: Manage Instances
1. Select an instance from dropdown
2. View instance details
3. Test connection
4. Delete instance (with confirmation)

## API Endpoints Used

### Get All Instances
```
GET /api/instances
Response: {
    "instances": [
        {
            "instance_id": "uuid",
            "product_id": "string",
            "version": "string",
            "registered_at": "timestamp"
        }
    ]
}
```

### Register Instance
```
POST /api/instance/register
Request: {
    "product_id": "string",
    "version": "string",
    "lcc_url": "string"
}
Response: {
    "success": bool,
    "instance_id": "uuid",
    "product_id": "string",
    "version": "string",
    "registered_at": "timestamp"
}
```

### Test Instance
```
POST /api/instance/test
Request: {
    "product_id": "string",
    "instance_id": "uuid"
}
Response: {
    "success": bool,
    "status": "string",
    "details": {}
}
```

### Delete Instance
```
POST /api/instance/clear
Request: {
    "product_id": "string",
    "instance_id": "uuid"
}
Response: {
    "success": bool
}
```

## Technical Details

### Instance Isolation
- Each instance has a unique instance ID
- Independent quota tracking per instance
- Separate authentication keys per instance
- Isolated state management

### Concurrent Registration
- Instances register sequentially (safer)
- Progress updates in real-time
- Error handling per instance
- Partial success supported

### Data Persistence
- Instances stored in backend (session or database)
- Survives page refresh (if backend supports)
- Can be cleared individually or in batch

## Limitations & Future Enhancements

### Current Limitations
1. Instances stored in memory (lost on server restart)
2. No persistent storage between sessions
3. No export/import of instance configurations

### Planned Enhancements
1. Persistent storage with database
2. Instance configuration export (JSON/YAML)
3. Bulk import of instances
4. Instance grouping/tagging
5. Instance usage statistics
6. Automated backup of instances

## Migration from Single-Instance Setup

### Old Flow (Single Instance)
```
1. Select one product
2. Configure SDK
3. Generate/load keys
4. Register (replaces previous)
5. Test single instance
```

### New Flow (Multi-Instance)
```
1. Select multiple products (checkboxes)
2. Set version (applies to all)
3. Register all at once
4. View dashboard with all instances
5. Manage and test each individually
```

## Example: Multi-Version Testing

### Scenario
Test if Product X behaves correctly when upgraded from v1.0.0 to v2.0.0

### Steps
1. Navigate to Setup page
2. Select "Product X" and set version to "1.0.0"
3. Click "Register Selected Products"
4. Instance created: `product-x-v1.0.0` (Instance ID: `abc-123`)
5. Select "Product X" again and set version to "2.0.0"
6. Click "Register Selected Products"
7. Instance created: `product-x-v2.0.0` (Instance ID: `abc-456`)
8. In management panel, switch between instances and test each

### Quota Tracking
- Each instance tracks quota independently
- v1.0.0 instance uses 100 units → 900 remaining
- v2.0.0 instance uses 100 units → 900 remaining (independent bucket)
- Both instances can be tested simultaneously
