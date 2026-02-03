import SwiftUI

struct LoginView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @State private var showRegister = false
    @State private var selectedRole: UserRole = .fan

    var body: some View {
        VStack(spacing: AppSpacing.xl) {
            Spacer()

            VStack(spacing: AppSpacing.md) {
                Image(systemName: "star.fill")
                    .font(.system(size: 60))
                    .foregroundStyle(Color.srGradientPrimary)

                GradientText("StarRink", font: .srHeading1)

                Text("Хоккейная платформа")
                    .font(.srBody)
                    .foregroundColor(.srTextSecondary)
            }

            #if DEBUG
            debugRoleButtons
            #endif

            VStack(spacing: AppSpacing.md) {
                SRTextField(placeholder: "Email", text: $authViewModel.email)
                SRTextField(placeholder: "Пароль", text: $authViewModel.password, isSecure: true)

                if let error = authViewModel.errorMessage {
                    Text(error)
                        .font(.srCaption)
                        .foregroundColor(.srError)
                }

                PrimaryButton("Войти", isLoading: authViewModel.isLoading) {
                    Task { await authViewModel.login() }
                }
            }
            .glassCard()

            Spacer()

            VStack(spacing: AppSpacing.xs) {
                Text("Нет аккаунта?")
                    .font(.srBody)
                    .foregroundColor(.srTextSecondary)

                Button("Создать аккаунт") {
                    showRegister = true
                }
                .font(.srBodyMedium)
                .foregroundColor(.srCyan)
            }
            .padding(.bottom, AppSpacing.xl)
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
        .sheet(isPresented: $showRegister) {
            RegisterView()
        }
    }

    #if DEBUG
    private var debugRoleButtons: some View {
        VStack(spacing: AppSpacing.sm) {
            Text("Быстрый вход (DEBUG)")
                .font(.srCaption)
                .foregroundColor(.srTextMuted)

            HStack(spacing: 0) {
                ForEach(UserRole.allCases, id: \.self) { role in
                    Button {
                        withAnimation(.spring(response: 0.3, dampingFraction: 0.75)) {
                            selectedRole = role
                        }
                        let impact = UIImpactFeedbackGenerator(style: .light)
                        impact.impactOccurred()
                    } label: {
                        VStack(spacing: 4) {
                            Image(systemName: role.icon)
                                .font(.system(size: 18))
                            Text(role.displayName)
                                .font(.system(size: 10, weight: .medium))
                        }
                        .foregroundColor(selectedRole == role ? .srCyan : .srTextMuted)
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, AppSpacing.sm)
                        .background(
                            RoundedRectangle(cornerRadius: 12)
                                .fill(selectedRole == role ? Color.srCyan.opacity(0.15) : Color.clear)
                        )
                        .overlay(
                            RoundedRectangle(cornerRadius: 12)
                                .stroke(
                                    selectedRole == role ? Color.srCyan.opacity(0.4) : Color.clear,
                                    lineWidth: 1
                                )
                        )
                        .scaleEffect(selectedRole == role ? 1.05 : 1.0)
                    }
                    .buttonStyle(.plain)
                }
            }
            .glassCard(padding: AppSpacing.sm)

            Button {
                let impact = UIImpactFeedbackGenerator(style: .medium)
                impact.impactOccurred()
                authViewModel.mockLogin(role: selectedRole)
            } label: {
                Text("Войти как \(selectedRole.displayName)")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundColor(.white)
                    .frame(maxWidth: .infinity)
                    .frame(height: 40)
                    .background(
                        RoundedRectangle(cornerRadius: 12)
                            .fill(
                                LinearGradient(
                                    colors: [.srCyan, .srPurple],
                                    startPoint: .leading,
                                    endPoint: .trailing
                                )
                            )
                    )
            }
        }
    }
    #endif
}
