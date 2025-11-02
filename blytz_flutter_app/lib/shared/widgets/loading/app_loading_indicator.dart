import 'package:flutter/material.dart';
import 'package:getwidget/getwidget.dart';

class AppLoadingIndicator extends StatelessWidget {
  final Color? color;
  final double? size;
  final String? message;

  const AppLoadingIndicator({
    super.key,
    this.color,
    this.size,
    this.message,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        GFLoader(
          loaderColorOne: color ?? Theme.of(context).primaryColor,
          loaderColorTwo: color ?? Theme.of(context).primaryColor,
          loaderColorThree: color ?? Theme.of(context).primaryColor,
          size: size ?? GFSize.MEDIUM,
        ),
        if (message != null) ...[
          const SizedBox(height: 16),
          Text(
            message!,
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ],
    );
  }
}

class AppLoadingOverlay extends StatelessWidget {
  final Widget child;
  final bool isLoading;
  final String? message;

  const AppLoadingOverlay({
    super.key,
    required this.child,
    required this.isLoading,
    this.message,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        child,
        if (isLoading)
          Container(
            color: Colors.black.withOpacity(0.5),
            child: Center(
              child: AppLoadingIndicator(
                color: Colors.white,
                message: message,
              ),
            ),
          ),
      ],
    );
  }
}

class AppLoadingButton extends StatelessWidget {
  final bool isLoading;
  final Widget child;
  final VoidCallback? onPressed;
  final ButtonStyle? style;

  const AppLoadingButton({
    super.key,
    required this.isLoading,
    required this.child,
    this.onPressed,
    this.style,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: isLoading ? null : onPressed,
      style: style,
      child: isLoading
          ? const SizedBox(
              height: 20,
              width: 20,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
              ),
            )
          : child,
    );
  }
}