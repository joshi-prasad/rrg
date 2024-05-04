import matplotlib.pyplot as plt
import matplotlib as mpl
import json
from matplotlib.colors import hsv_to_rgb
from cycler import cycler

def plot_multi_rsi_vs_rs(data, title="RSI vs RS Chart"):
  """
  Plots RSI values on the vertical axis and RS values on the horizontal axis
  for multiple stocks on separate lines with labels.

  Args:
      data: A dictionary where keys are stock names (strings) and values are
             tuples containing RSI and RS values as lists (e.g., {"StockA": (rsi_a, rs_a), ...})
      title: The title of the chart (optional).
  """

  # Set plot limits (0-100 for both axes)
  plt.xlim(0, 100)
  plt.ylim(35, 100)

  # Loop through each stock data and plot
  # 1000 distinct colors:
  colors = [hsv_to_rgb([(i * 0.618033988749895) % 1.0, 1, 1])
            for i in range(1000)]
  plt.rc('axes', prop_cycle=(cycler('color', colors)))
  for stock_name, rs_rsi_json in data.items():
    # Scatter plot and line for each stock
    rs_values = rs_rsi_json["rs"]
    rsi_values = rs_rsi_json["rsi"]
    plt.scatter(rs_values, rsi_values, marker='o', alpha=0.7, label=stock_name)
    plt.plot(rs_values, rsi_values, '-b', alpha=0.7, label=stock_name)

  # Add labels and title
  plt.xlabel('RS Value (0-100)')
  plt.ylabel('RSI Value (1-100)')
  plt.title(title)

  # Add legend
  plt.legend()

  # Display the plot
  plt.grid(visible=False)
  plt.show()


with open("data.json") as fh:
  json_data = json.load(fh)
  plot_multi_rsi_vs_rs(json_data)
