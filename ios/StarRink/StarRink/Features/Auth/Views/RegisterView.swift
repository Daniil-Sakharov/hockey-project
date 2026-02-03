import SwiftUI

struct RegisterView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack {
            Color.srBackground.ignoresSafeArea()

            VStack(spacing: AppSpacing.xl) {
                HStack {
                    Button { dismiss() } label: {
                        Image(systemName: "xmark")
                            .font(.title3)
                            .foregroundColor(.srTextSecondary)
                    }
                    Spacer()
                }
                .padding(.top)

                VStack(spacing: AppSpacing.md) {
                    GradientText("Регистрация", font: .srHeading2)
                    Text("Создайте аккаунт StarRink")
                        .font(.srBody)
                        .foregroundColor(.srTextSecondary)
                }

                VStack(spacing: AppSpacing.md) {
                    SRTextField(placeholder: "Имя", text: $authViewModel.name)
                    SRTextField(placeholder: "Email", text: $authViewModel.email)
                    SRTextField(placeholder: "Пароль", text: $authViewModel.password, isSecure: true)
                    SRTextField(placeholder: "Подтвердите пароль", text: $authViewModel.confirmPassword, isSecure: true)

                    if let error = authViewModel.errorMessage {
                        Text(error)
                            .font(.srCaption)
                            .foregroundColor(.srError)
                    }

                    PrimaryButton("Создать аккаунт", isLoading: authViewModel.isLoading) {
                        Task {
                            await authViewModel.register()
                            if authViewModel.isAuthenticated { dismiss() }
                        }
                    }
                }
                .glassCard()

                Spacer()
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
        }
    }
}
