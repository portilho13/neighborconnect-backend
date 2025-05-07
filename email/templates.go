package email

import (
	"fmt"

	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
)

func CreateRewardEmailTemplate(event models.Community_Event) string {
	return fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Your Event Reward</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					color: #333333;
					max-width: 600px;
					margin: 0 auto;
					padding: 0;
					background-color: #f9fafb;
				}
				.container {
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
					margin: 20px;
					border: 1px solid #e5e7eb;
				}
				.header {
					background-color: #3F3D56;
					color: #ffffff;
					padding: 24px;
					text-align: center;
				}
				.content {
					background-color: #ffffff;
					padding: 24px;
					text-align: center;
				}
				.footer {
					background-color: #f6f6f6;
					padding: 16px 24px;
					text-align: center;
					font-size: 12px;
					color: #666666;
					border-top: 1px solid #e5e7eb;
				}
				h1 {
					margin: 0;
					font-size: 22px;
					font-weight: 700;
				}
				h2 {
					color: #3F3D56;
					font-size: 20px;
					margin-top: 0;
				}
				.reward-code-container {
					background: linear-gradient(to right, rgba(63, 61, 86, 0.1), rgba(108, 106, 138, 0.1));
					border-radius: 8px;
					padding: 24px;
					margin: 24px 0;
					border: 1px dashed #3F3D56;
				}
				.reward-code {
					font-family: 'Courier New', monospace;
					font-size: 28px;
					font-weight: bold;
					letter-spacing: 2px;
					color: #3F3D56;
					padding: 12px 24px;
					background-color: #ffffff;
					border-radius: 6px;
					display: inline-block;
					margin: 10px 0;
					border: 1px solid #e5e7eb;
				}
				.instructions {
					background-color: #f9fafb;
					border: 1px solid #e5e7eb;
					border-radius: 8px;
					padding: 16px;
					margin: 16px 0;
					text-align: left;
				}
				.instructions ol {
					margin: 0;
					padding-left: 24px;
				}
				.instructions li {
					margin-bottom: 8px;
				}
				.button {
					display: inline-block;
					background-color: #3F3D56;
					color: #ffffff;
					text-decoration: none;
					padding: 12px 24px;
					border-radius: 8px;
					font-weight: 500;
					margin-top: 16px;
				}
				.expiry {
					color: #ef4444;
					font-weight: 500;
					margin-top: 16px;
				}
				.social-links {
					margin-top: 16px;
				}
				.social-links a {
					display: inline-block;
					margin: 0 8px;
					color: #666666;
					text-decoration: none;
				}
				@media only screen and (max-width: 600px) {
					.container {
						margin: 10px;
					}
					.header, .content, .footer {
						padding: 16px;
					}
					.reward-code {
						font-size: 22px;
						padding: 10px 16px;
					}
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Congratulations! You Earned a Reward</h1>
				</div>
				<div class="content">
					<h2>Thank you for participating in our event</h2>
					<p>We’re happy to reward you for your participation in the NeighboorConnect community event.</p>
					
					<div class="reward-code-container">
						<p>Your reward code is:</p>
						<div class="reward-code">%s</div>
						<p>Use this code to redeem your reward</p>
					</div>
					
					<div class="instructions">
						<p><strong>How to use your code:</strong></p>
						<ol>
							<li>Log in to your NeighboorConnect account</li>
							<li>Go to the dashboard</li>
							<li>Click “Redeem Code”</li>
							<li>Enter the code above and click “Apply”</li>
						</ol>
					</div>
					
					<p class="expiry">This code expires in: 5 days</p>
					
					<a href="#" class="button">Redeem Now</a>
					
					<p style="margin-top: 24px;">Enjoy your reward!<br>NeighboorConnect Team</p>
				</div>
				<div class="footer">
					<p>This is an automated message. Please do not reply to this email.</p>
					<div class="social-links">
						<a href="#">Facebook</a> • 
						<a href="#">Instagram</a> • 
						<a href="#">Website</a>
					</div>
					<p>© 2024 NeighboorConnect. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>`, event.Code,
	)
}
