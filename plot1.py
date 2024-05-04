import matplotlib.pyplot as plt
import matplotlib as mpl
import json
from matplotlib.colors import hsv_to_rgb

def plot_multi_rsi_vs_rs(data, title="RSI vs RS Chart"):
  """
  Plots RSI values on the vertical axis and RS values on the horizontal axis
  for multiple stocks on separate lines, using arrows to connect scatter plot points
  and printing the stock name once per line. All elements for a stock use the same color.

  Args:
    data: A dictionary where keys are stock names (strings) and values are
      tuples containing RSI and RS values as lists (e.g., {"StockA": (rsi_a, rs_a), ...})
    title: The title of the chart (optional).
  """

  # Initialize plot limits with large values
  min_rs, max_rs = float('inf'), float('-inf')
  min_rsi, max_rsi = float('inf'), float('-inf')

  # Loop through each stock data to find minimum and maximum values
  for rs_rsi_json in data.values():
    rs_values = rs_rsi_json["rs"][-6:]
    rsi_values = rs_rsi_json["rsi_ema"][-6:]
    min_rs = min(min_rs, min(rs_values))
    max_rs = max(max_rs, max(rs_values))
    min_rsi = min(min_rsi, min(rsi_values))
    max_rsi = max(max_rsi, max(rsi_values))

  # Set plot limits based on the minimum and maximum values
  plt.xlim(min_rs - 1, max_rs + 1)
  plt.ylim(min_rsi - 1, max_rsi + 1)
  # Set logarithmic scale for x-axis
  # plt.xscale('log')
  # plt.yscale('log')

  # Loop through each stock data and plot
  color_iter = iter(plt.cm.get_cmap('tab20').colors)  # Colormap iterator for multiple stocks
  for stock_name, rs_rsi_json in data.items():
    # if stock_name not in ["INDEXNSE:NIFTY_METAL", "INDEXNSE:NIFTY_PHARMA"]:
    #   continue
    # Choose a color for the stock
    color = next(color_iter)

    # Extract data and format for plotting
    rs_values = rs_rsi_json["rs"][-6:]
    rsi_values = rs_rsi_json["rsi_ema"][-6:]
    print(stock_name, rsi_values)
    # Scatter plot, arrows, and text (all with the same color)
    plt.scatter(rs_values, rsi_values, marker='o', alpha=0.7, color=color, label=stock_name)
    for i in range(len(rs_values) - 1):
      dx = rs_values[i + 1] - rs_values[i]
      dy = rsi_values[i + 1] - rsi_values[i]
      arrow = plt.arrow(rs_values[i], rsi_values[i], dx, dy, head_width=0.3, head_length=0.5, color=color)
    plt.text(rs_values[0] + 0.5, rsi_values[0] + 0.2, stock_name, ha="center", va="center", fontsize=8, color=color)

  # Add labels and title
  plt.xlabel('RS Value')
  plt.ylabel('RSI Value')
  plt.title(title)

  # Adjust margins to reduce whitespace
  plt.subplots_adjust(left=0.05, right=0.95, top=0.95, bottom=0.05)

  # Display the plot
  plt.grid(visible=False)
  plt.show()



with open("data.json") as fh:
  json_data = json.load(fh)
  plot_multi_rsi_vs_rs(json_data)
