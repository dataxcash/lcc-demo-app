# Week 4 Quick Start Guide

## 🚀 Testing the Instance Setup Page

### Prerequisites

1. **LCC Server Running**
   ```bash
   # Make sure your LCC server is running on localhost:7086
   # Or adjust the URL in the demo app
   ```

2. **Demo App Built**
   ```bash
   cd /home/fila/jqdDev_2025/lcc-demo-app
   go build -o lcc-demo-app cmd/webdemo/main.go
   ```

### Step-by-Step Testing

#### 1. Start the Demo App

```bash
./lcc-demo-app
# Server will start on http://localhost:9144
```

#### 2. Open in Browser

Navigate to: http://localhost:9144

#### 3. Configure LCC Connection (Week 1)

1. Click on "Welcome" or navigate to the first page
2. Enter LCC Server URL: `http://localhost:7086`
3. Click "Test Connection"
4. Click "Save & Continue"

#### 4. Navigate to Instance Setup Page

1. Click on "Setup" in the navigation
2. Or use the "Next" button to progress through pages
3. URL should be: `http://localhost:9144/#setup`

#### 5. Test Instance Registration Flow

##### Step 1: Select Product
1. Click "Refresh Products" to load from LCC server
2. Select a product from dropdown (e.g., "data-insight-pro")
3. Verify product details appear below

##### Step 2: Configure SDK Client
1. Product ID auto-fills from selection
2. Set Product Version (default: 1.0.0)
3. LCC URL should auto-load from config

##### Step 3: Authentication Keys
1. Click "Generate New Key Pair" (optional)
2. Or click "Use Saved Keys" to use keystore
3. Keys are handled automatically during registration

##### Step 4: Register Instance
1. Click "Register Instance" button
2. Watch for "Registering..." loading state
3. Success message should appear
4. Instance Status card should appear with:
   - Instance ID (unique identifier)
   - Product ID
   - Active status
   - Registered timestamp

##### Test Instance Connection
1. Select a feature from dropdown (e.g., "ml_analytics")
2. Click "Test CheckFeature()"
3. JSON result should display below
4. Border color indicates success (green) or failure (red)

##### Clear Instance (Optional)
1. Click "Clear Instance" button
2. Instance status card should hide
3. Register button should reappear
4. Can re-register the same product

---

## 🧪 API Testing (curl)

### Test Registration
```bash
curl -X POST http://localhost:9144/api/instance/register \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "data-insight-pro",
    "version": "1.0.0",
    "lcc_url": "http://localhost:7086"
  }'
```

Expected response:
```json
{
  "success": true,
  "instance_id": "inst-abc123...",
  "product_id": "data-insight-pro",
  "version": "1.0.0",
  "registered_at": "2025-01-21T12:00:00Z",
  "message": "Instance registered successfully"
}
```

### Test Status Check
```bash
curl "http://localhost:9144/api/instance/status?product_id=data-insight-pro"
```

### Test Feature Check
```bash
curl -X POST http://localhost:9144/api/instance/test \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "data-insight-pro",
    "feature_id": "ml_analytics"
  }'
```

### Clear Instance
```bash
curl -X POST http://localhost:9144/api/instance/clear \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "data-insight-pro"
  }'
```

### Generate Keys
```bash
curl -X POST http://localhost:9144/api/instance/generate-keys
```

---

## ✅ Verification Checklist

### UI Functionality
- [ ] Page loads without console errors
- [ ] Products dropdown populates
- [ ] Product selection shows details
- [ ] Config fields auto-populate correctly
- [ ] Register button disabled until product selected
- [ ] Registration succeeds with valid data
- [ ] Success message displays
- [ ] Instance ID shows in status card
- [ ] Feature test returns JSON
- [ ] Clear instance resets UI
- [ ] Can re-register after clear

### API Functionality
- [ ] POST /api/instance/register works
- [ ] GET /api/instance/status works
- [ ] POST /api/instance/test works
- [ ] POST /api/instance/clear works
- [ ] POST /api/instance/generate-keys works

### Error Handling
- [ ] Registration fails gracefully with invalid product
- [ ] Error messages are user-friendly
- [ ] No product selected prevents registration
- [ ] Missing LCC URL shows error
- [ ] LCC server unreachable shows error

### Integration
- [ ] Uses config from Week 1 (Welcome page)
- [ ] Products list matches Week 2 (Tiers page)
- [ ] Instance works with Week 3 (Limits page)

---

## 🐛 Troubleshooting

### Issue: "Failed to load products"

**Cause:** LCC server not configured or unreachable

**Solution:**
1. Go to Welcome page (#welcome)
2. Configure LCC URL
3. Test connection
4. Save configuration
5. Return to Setup page

### Issue: "Registration failed: product not found"

**Cause:** Product doesn't exist on LCC server

**Solution:**
1. Click "Refresh Products" on Setup page
2. Select from available products
3. Verify product exists on LCC server

### Issue: "Instance already registered"

**Cause:** Product already has active client

**Solution:**
1. Click "Clear Instance" button
2. Or restart demo app to clear in-memory state

### Issue: Keys not persisting

**Cause:** Keystore directory permissions

**Solution:**
```bash
# Check keystore directory
ls -la ~/.lcc-demo/keys/

# Fix permissions if needed
chmod 700 ~/.lcc-demo/keys/
```

---

## 📊 Expected Behavior

### First Registration
1. ⏱️ Registration takes 1-3 seconds
2. ✅ Success message appears
3. 🆔 Instance ID displays (format: `inst-...`)
4. 📦 Instance stored in server memory
5. 🔑 Keys saved to `~/.lcc-demo/keys/{product_id}.json`

### Subsequent Registrations (same product)
1. 🔑 Loads existing keys from keystore
2. 🆔 New instance ID generated
3. ♻️ Old client replaced in memory

### Feature Testing
1. 📡 Calls CheckFeature() on registered instance
2. 🎯 Returns enabled/disabled status
3. 📝 Shows reason and details
4. 🎨 Visual feedback with color

---

## 📁 File Locations

### Backend
- Handler: `internal/web/instance_handler.go` (345 lines)
- Routes: `internal/web/server.go` (added 5 routes)

### Frontend
- Page: `static/js/pages/setup.js` (473 lines)

### Data Storage
- Config: `~/.lcc-demo/config.json`
- Keys: `~/.lcc-demo/keys/{product_id}.json`

---

## 🎓 Learning Points to Verify

Users should understand after testing:

1. **SDK Client = Interface to LCC**
   - One client per product
   - Manages all SDK operations
   - Requires configuration

2. **Registration = Establishing Trust**
   - Sends public key to LCC
   - Receives instance ID
   - Validates connection

3. **Instance ID = Unique Identifier**
   - Generated by LCC server
   - Used for tracking and auditing
   - Persists across app restarts

4. **Key Pairs = Security**
   - RSA public/private keys
   - Stored securely in keystore
   - Reusable across registrations

5. **Product ID = License Context**
   - Which product to manage
   - Determines available features
   - Required for all operations

---

## 🎯 Success Criteria

✅ Registration completes successfully  
✅ Instance ID displays correctly  
✅ Feature tests return valid results  
✅ Error messages are helpful  
✅ UI is intuitive and responsive  
✅ Code examples are educational  
✅ Integrates with previous pages  

---

## 📞 Support

If you encounter issues:

1. Check browser console for errors
2. Verify LCC server is running
3. Check server logs for backend errors
4. Review API responses with curl
5. Restart demo app to clear state

---

**Happy Testing! 🚀**
