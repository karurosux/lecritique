# Kyooar Tasks

## 🚧 Active Development

### Questionnaire Builder

**Status**: 80% complete

**Remaining**:

- [ ] Drag-and-drop question ordering
- [ ] Question preview
- [ ] Question templates
- [ ] Bulk operations

## 🐛 Known Issues

- Question ordering doesn't persist
- Settings page incomplete
- Email templates need work
- QR code API endpoint issue (FindByCode)

## 📋 Next Priority Tasks

### Essential Features

1. **Customer feedback submission flow**
   - Feedback form page with dynamic questionnaire loading
   - Question types: rating, scale, single choice, yes/no, text
   - Success confirmation page
   - Anonymous submission support

2. **Organization owner dashboard**
   - Statistics overview cards
   - Recent feedback display
   - Quick actions menu
   - Organization switcher

3. **Organization and product management**
   - Organization list/grid view
   - Add/edit organization form
   - Product management interface
   - Category management

4. **QR code generation**
   - Generate new QR codes with types
   - Download QR codes (PNG/SVG)
   - Print-friendly QR code sheets
   - QR code analytics

5. **Feedback analytics**
   - Rating trends charts
   - Feedback volume over time
   - Product performance comparison
   - Export functionality

## 💡 Backlog Ideas

- AI response analysis
- WhatsApp integration
- Multi-language support
- PWA with offline mode
- Email notifications (low rating alerts, summaries)
- Advanced analytics (sentiment analysis, word clouds)

## 🔧 Technical Debt

- [ ] Add comprehensive error boundaries
- [ ] Implement proper logging
- [ ] Add unit tests for components
- [ ] Add E2E tests with Playwright
- [ ] Optimize bundle size
- [ ] Implement proper caching strategy

## 📝 Development Notes

- **Test User**: admin@kyooar.com / admin123
- **Commands**: `make seed` or `make seed-force`
- **Remember**: Questionnaires are product-specific (core feature!)
- **Always**: Check subscription limits before operations
- **UI**: Use Svelte 5 runes ($state, $derived, $effect)

## Last Updated

2025-07-17
