# Blytz Flutter App

A production-ready Flutter mobile application for the Blytz Live Auction Platform.

## 🚀 Features

### ✅ Completed Features
- **Project Structure**: Clean architecture with feature-based organization
- **Core Services**: API client, authentication, secure storage, local database
- **State Management**: Riverpod with comprehensive providers
- **Navigation**: Go Router with type-safe routing and guards
- **Authentication**: Complete login/register/forgot password flow
- **UI Framework**: Material Design 3 with custom theming
- **Form Validation**: Comprehensive validation utilities
- **Error Handling**: Centralized error management
- **Storage**: Secure storage, local database, and preferences
- **Networking**: Type-safe API client with interceptors

### 🚧 In Progress
- Auction listing and detail screens
- Real-time bidding with WebSocket
- LiveKit video streaming integration

### 📋 Planned Features
- Real-time chat functionality
- Payment processing integration
- Push notifications
- Internationalization
- Analytics and monitoring
- Offline support

## 🏗️ Architecture

### Project Structure
```
lib/
├── core/                    # Core functionality
│   ├── constants/          # App constants
│   ├── errors/             # Error handling
│   ├── network/            # API client and networking
│   ├── services/           # Core services
│   ├── storage/            # Storage solutions
│   ├── utils/              # Utilities and helpers
│   └── providers/          # Global state providers
├── features/               # Feature modules
│   ├── auth/              # Authentication
│   ├── auctions/          # Auction functionality
│   ├── chat/              # Chat functionality
│   ├── payments/          # Payment processing
│   ├── profile/           # User profile
│   └── onboarding/        # Onboarding flow
├── shared/                # Shared components
│   ├── widgets/           # Reusable widgets
│   ├── themes/            # App theming
│   └── pages/             # Shared pages
└── main.dart              # App entry point
```

### Technology Stack
- **Framework**: Flutter 3.16.0+
- **State Management**: Riverpod 2.4.9
- **Navigation**: Go Router 12.1.3
- **Networking**: Dio 5.4.0 + Retrofit 4.0.3
- **Storage**: Hive 2.2.3 + Flutter Secure Storage 9.0.0
- **UI Components**: GetWidget 4.0.0
- **Real-time**: WebSocket + LiveKit
- **Validation**: Form Validator 2.1.1

## 📦 Dependencies

### Core Dependencies
- `flutter_riverpod`: State management
- `go_router`: Navigation and routing
- `dio`: HTTP client
- `retrofit`: Type-safe API generation
- `hive`: Local database
- `flutter_secure_storage`: Secure storage
- `getwidget`: UI components
- `form_validator`: Form validation

### Development Dependencies
- `build_runner`: Code generation
- `retrofit_generator`: API code generation
- `json_serializable`: JSON serialization
- `hive_generator`: Hive code generation
- `mockito`: Testing framework

## 🔧 Setup

### Prerequisites
- Flutter SDK 3.16.0 or higher
- Dart SDK compatible with Flutter version
- Android Studio / VS Code with Flutter extensions

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd blytz_flutter_app
   ```

2. Install dependencies:
   ```bash
   flutter pub get
   ```

3. Generate code:
   ```bash
   flutter packages pub run build_runner build
   ```

4. Run the app:
   ```bash
   flutter run
   ```

## 🔐 Environment Configuration

Create a `.env` file in the root directory:

```env
API_BASE_URL=https://api.blytz.app
WS_URL=wss://api.blytz.app
LIVEKIT_URL=wss://livekit.blytz.app
SENTRY_DSN=your_sentry_dsn
```

## 🧪 Testing

### Run Tests
```bash
# Unit tests
flutter test

# Widget tests
flutter test test/widget/

# Integration tests
flutter test integration_test/
```

### Test Coverage
```bash
flutter test --coverage
genhtml coverage/lcov.info -o coverage/html
```

## 🚀 Build & Deployment

### Android
```bash
# Debug APK
flutter build apk --debug

# Release APK
flutter build apk --release

# App Bundle
flutter build appbundle --release
```

### iOS
```bash
# Debug build
flutter build ios --debug

# Release build
flutter build ios --release
```

## 📱 Screenshots

### Authentication Flow
- Login screen with form validation
- Registration with multi-step form
- Forgot password flow

### Main App
- Home screen with navigation
- Auction listings
- User profile

## 🔄 State Management

### Providers
- `authStateProvider`: Authentication state
- `apiClientProvider`: API client instance
- `secureStorageProvider`: Secure storage
- `localDatabaseProvider`: Local database
- `themeProvider`: Theme management
- `languageProvider`: Language preferences

### Data Flow
1. UI triggers action via provider
2. Provider calls service layer
3. Service interacts with API/local storage
4. State updates and UI rebuilds

## 🔗 API Integration

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Token refresh
- `POST /api/v1/auth/logout` - User logout

### Auction Endpoints
- `GET /api/v1/auctions` - List auctions
- `GET /api/v1/auctions/:id` - Get auction details
- `POST /api/v1/auctions/:id/bids` - Place bid

## 🎨 UI Components

### Custom Widgets
- `AppTextField`: Custom text field with validation
- `AppLoadingButton`: Loading state button
- `AppLoadingIndicator`: Custom loading spinner
- `AuctionCard`: Auction listing card
- `BidButton`: Bid placement button

### Theme
- Light and dark theme support
- Custom color scheme
- Material Design 3 components

## 🔒 Security

### Token Management
- JWT tokens stored in secure storage
- Automatic token refresh
- Token expiration handling

### Data Protection
- Sensitive data encrypted at rest
- HTTPS for all API calls
- Certificate pinning (future)

## 📈 Performance

### Optimizations
- Lazy loading for large lists
- Image caching with `cached_network_image`
- Efficient state management with Riverpod
- Code splitting for reduced bundle size

### Monitoring
- Performance tracking with Firebase Performance
- Error reporting with Sentry
- Analytics with Firebase Analytics

## 🌍 Internationalization

### Supported Languages
- English (en)
- Spanish (es)
- French (fr)
- Chinese (zh)
- Japanese (ja)

### Implementation
- `easy_localization` for i18n
- Dynamic language switching
- RTL language support

## 📝 Contributing

### Development Workflow
1. Create feature branch from `develop`
2. Implement changes with tests
3. Run linting and tests
4. Submit pull request

### Code Style
- Follow `very_good_analysis` linting rules
- Use meaningful variable names
- Document public APIs
- Write tests for new features

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check documentation

## 🗺️ Roadmap

### Phase 1 (Current)
- ✅ Core architecture
- ✅ Authentication
- ✅ Basic UI components
- 🚧 Auction screens

### Phase 2 (Next)
- Real-time bidding
- Live streaming
- Chat functionality
- Payment integration

### Phase 3 (Future)
- Advanced analytics
- Machine learning features
- Enhanced security
- Performance optimizations